package kubernetes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

type Kubernetes struct {
	cfg                         *rest.Config
	clientSet                   *kubernetes.Clientset
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
	return &Kubernetes{
		cfg:                         cfg,
		clientSet:                   clientSet,
		discoveryClient:             discoveryClient,
		deferredDiscoveryRESTMapper: restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient)),
		dynamicClient:               dynamicClient,
	}, nil
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
	//cfg, err := pathOptions.GetStartingConfig()
	//if err != nil {
	//	return nil, err
	//}
	//if err = clientcmdapi.MinifyConfig(cfg); err != nil {
	//	return nil, err
	//}
	//if err = clientcmdapi.FlattenConfig(cfg); err != nil {
	//	return nil, err
	//}
	//return cfg, nil
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: pathOptions.GetDefaultFilename()},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})
}

func resolveClientConfig() (*rest.Config, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err == nil && inClusterConfig != nil {
		return inClusterConfig, nil
	}
	return resolveConfig().ClientConfig()
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
