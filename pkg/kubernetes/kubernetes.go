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
	cfg, err := resolveClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	scheme := runtime.NewScheme()
	if err = v1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return &Kubernetes{
		cfg:                         cfg,
		kubeConfigFiles:             resolveConfig().ConfigAccess().GetLoadingPrecedence(),
		scheme:                      scheme,
		parameterCodec:              runtime.NewParameterCodec(scheme),
		clientSet:                   clientSet,
		discoveryClient:             discoveryClient,
		deferredDiscoveryRESTMapper: restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient)),
		dynamicClient:               dynamicClient,
	}, nil
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
