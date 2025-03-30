package kubernetes

import (
	"github.com/fsnotify/fsnotify"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

// InClusterConfig is a variable that holds the function to get the in-cluster config
// Exposed for testing
var InClusterConfig = func() (*rest.Config, error) {
	// TODO use kubernetes.default.svc instead of resolved server
	// Currently running into: `http: server gave HTTP response to HTTPS client`
	inClusterConfig, err := rest.InClusterConfig()
	if inClusterConfig != nil {
		inClusterConfig.Host = "https://kubernetes.default.svc"
	}
	return inClusterConfig, err
}

type CloseWatchKubeConfig func() error

type Kubernetes struct {
	cfg                         *rest.Config
	kubeConfigFiles             []string
	CloseWatchKubeConfig        CloseWatchKubeConfig
	scheme                      *runtime.Scheme
	parameterCodec              runtime.ParameterCodec
	clientSet                   kubernetes.Interface
	discoveryClient             *discovery.DiscoveryClient
	deferredDiscoveryRESTMapper *restmapper.DeferredDiscoveryRESTMapper
	dynamicClient               *dynamic.DynamicClient
}

func NewKubernetes() (*Kubernetes, error) {
	k8s := &Kubernetes{}
	var err error
	k8s.cfg, err = resolveClientConfig()
	if err != nil {
		return nil, err
	}
	k8s.kubeConfigFiles = resolveConfig().ConfigAccess().GetLoadingPrecedence()
	k8s.clientSet, err = kubernetes.NewForConfig(k8s.cfg)
	if err != nil {
		return nil, err
	}
	k8s.discoveryClient, err = discovery.NewDiscoveryClientForConfig(k8s.cfg)
	if err != nil {
		return nil, err
	}
	k8s.deferredDiscoveryRESTMapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(k8s.discoveryClient))
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

func (k *Kubernetes) WatchKubeConfig(onKubeConfigChange func() error) {
	if len(k.kubeConfigFiles) == 0 {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	for _, file := range k.kubeConfigFiles {
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
	if k.CloseWatchKubeConfig != nil {
		_ = k.CloseWatchKubeConfig()
	}
	k.CloseWatchKubeConfig = watcher.Close
}

func (k *Kubernetes) Close() {
	if k.CloseWatchKubeConfig != nil {
		_ = k.CloseWatchKubeConfig()
	}
}

func marshal(v any) (string, error) {
	switch t := v.(type) {
	case []unstructured.Unstructured:
		for i := range t {
			t[i].SetManagedFields(nil)
		}
	case []*unstructured.Unstructured:
		for i := range t {
			t[i].SetManagedFields(nil)
		}
	case unstructured.Unstructured:
		t.SetManagedFields(nil)
	case *unstructured.Unstructured:
		t.SetManagedFields(nil)
	}
	ret, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func resolveConfig() clientcmd.ClientConfig {
	pathOptions := clientcmd.NewDefaultPathOptions()
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: pathOptions.GetDefaultFilename()},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})
}

func resolveClientConfig() (*rest.Config, error) {
	inClusterConfig, err := InClusterConfig()
	if err == nil && inClusterConfig != nil {
		return inClusterConfig, nil
	}
	cfg, err := resolveConfig().ClientConfig()
	if cfg != nil && cfg.UserAgent == "" {
		cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	return cfg, err
}

func configuredNamespace() string {
	if ns, _, nsErr := resolveConfig().Namespace(); nsErr == nil {
		return ns
	}
	return ""
}

func namespaceOrDefault(namespace string) string {
	if namespace == "" {
		return configuredNamespace()
	}
	return namespace
}
