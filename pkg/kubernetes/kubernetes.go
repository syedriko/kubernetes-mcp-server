package kubernetes

import (
	"context"
	"strings"

	"github.com/fsnotify/fsnotify"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"

	"github.com/manusa/kubernetes-mcp-server/pkg/config"
	"github.com/manusa/kubernetes-mcp-server/pkg/helm"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const (
	AuthorizationHeader = "kubernetes-authorization"
)

type CloseWatchKubeConfig func() error

type Kubernetes struct {
	manager *Manager
}

type Manager struct {
	// Kubeconfig path override
	Kubeconfig                  string
	cfg                         *rest.Config
	clientCmdConfig             clientcmd.ClientConfig
	CloseWatchKubeConfig        CloseWatchKubeConfig
	scheme                      *runtime.Scheme
	parameterCodec              runtime.ParameterCodec
	clientSet                   kubernetes.Interface
	discoveryClient             discovery.CachedDiscoveryInterface
	deferredDiscoveryRESTMapper *restmapper.DeferredDiscoveryRESTMapper
	dynamicClient               *dynamic.DynamicClient

	StaticConfig *config.StaticConfig
}

func NewManager(kubeconfig string, config *config.StaticConfig) (*Manager, error) {
	k8s := &Manager{
		Kubeconfig:   kubeconfig,
		StaticConfig: config,
	}
	if err := resolveKubernetesConfigurations(k8s); err != nil {
		return nil, err
	}
	// TODO: Won't work because not all client-go clients use the shared context (e.g. discovery client uses context.TODO())
	//k8s.cfg.Wrap(func(original http.RoundTripper) http.RoundTripper {
	//	return &impersonateRoundTripper{original}
	//})
	var err error
	k8s.clientSet, err = kubernetes.NewForConfig(k8s.cfg)
	if err != nil {
		return nil, err
	}
	k8s.discoveryClient = memory.NewMemCacheClient(discovery.NewDiscoveryClient(k8s.clientSet.CoreV1().RESTClient()))
	k8s.deferredDiscoveryRESTMapper = restmapper.NewDeferredDiscoveryRESTMapper(k8s.discoveryClient)
	k8s.dynamicClient, err = dynamic.NewForConfig(k8s.cfg)
	if err != nil {
		return nil, err
	}
	k8s.scheme = runtime.NewScheme()
	if err = v1.AddToScheme(k8s.scheme); err != nil {
		return nil, err
	}
	k8s.parameterCodec = runtime.NewParameterCodec(k8s.scheme)
	return k8s, nil
}

func (m *Manager) WatchKubeConfig(onKubeConfigChange func() error) {
	if m.clientCmdConfig == nil {
		return
	}
	kubeConfigFiles := m.clientCmdConfig.ConfigAccess().GetLoadingPrecedence()
	if len(kubeConfigFiles) == 0 {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	for _, file := range kubeConfigFiles {
		_ = watcher.Add(file)
	}
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				_ = onKubeConfigChange()
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	if m.CloseWatchKubeConfig != nil {
		_ = m.CloseWatchKubeConfig()
	}
	m.CloseWatchKubeConfig = watcher.Close
}

func (m *Manager) Close() {
	if m.CloseWatchKubeConfig != nil {
		_ = m.CloseWatchKubeConfig()
	}
}

func (m *Manager) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return m.discoveryClient, nil
}

func (m *Manager) ToRESTMapper() (meta.RESTMapper, error) {
	return m.deferredDiscoveryRESTMapper, nil
}

func (m *Manager) Derived(ctx context.Context) *Kubernetes {
	authorization, ok := ctx.Value(AuthorizationHeader).(string)
	if !ok || !strings.HasPrefix(authorization, "Bearer ") {
		return &Kubernetes{manager: m}
	}
	klog.V(5).Infof("%s header found (Bearer), using provided bearer token", AuthorizationHeader)
	derivedCfg := rest.CopyConfig(m.cfg)
	derivedCfg.BearerToken = strings.TrimPrefix(authorization, "Bearer ")
	derivedCfg.BearerTokenFile = ""
	derivedCfg.Username = ""
	derivedCfg.Password = ""
	derivedCfg.AuthProvider = nil
	derivedCfg.AuthConfigPersister = nil
	derivedCfg.ExecProvider = nil
	derivedCfg.Impersonate = rest.ImpersonationConfig{}
	clientCmdApiConfig, err := m.clientCmdConfig.RawConfig()
	if err != nil {
		return &Kubernetes{manager: m}
	}
	clientCmdApiConfig.AuthInfos = make(map[string]*clientcmdapi.AuthInfo)
	derived := &Kubernetes{manager: &Manager{
		Kubeconfig:      m.Kubeconfig,
		clientCmdConfig: clientcmd.NewDefaultClientConfig(clientCmdApiConfig, nil),
		cfg:             derivedCfg,
		scheme:          m.scheme,
		parameterCodec:  m.parameterCodec,
	}}
	derived.manager.clientSet, err = kubernetes.NewForConfig(derived.manager.cfg)
	if err != nil {
		return &Kubernetes{manager: m}
	}
	derived.manager.discoveryClient = memory.NewMemCacheClient(discovery.NewDiscoveryClient(derived.manager.clientSet.CoreV1().RESTClient()))
	derived.manager.deferredDiscoveryRESTMapper = restmapper.NewDeferredDiscoveryRESTMapper(derived.manager.discoveryClient)
	derived.manager.dynamicClient, err = dynamic.NewForConfig(derived.manager.cfg)
	if err != nil {
		return &Kubernetes{manager: m}
	}
	return derived
}

func (k *Kubernetes) NewHelm() *helm.Helm {
	// This is a derived Kubernetes, so it already has the Helm initialized
	return helm.NewHelm(k.manager)
}
