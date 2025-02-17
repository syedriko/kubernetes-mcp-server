package kubernetes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

type Kubernetes struct {
	cfg                         *rest.Config
	deferredDiscoveryRESTMapper *restmapper.DeferredDiscoveryRESTMapper
}

func NewKubernetes() (*Kubernetes, error) {
	cfg, err := resolveClientConfig()
	if err != nil {
		return nil, err
	}
	return &Kubernetes{cfg: cfg}, nil
}

func marshal(v any) (string, error) {
	switch t := v.(type) {
	case []unstructured.Unstructured:
		for i := range t {
			t[i].SetManagedFields(nil)
		}
	case unstructured.Unstructured:
		t.SetManagedFields(nil)
	}
	ret, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func resolveClientConfig() (*rest.Config, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err == nil && inClusterConfig != nil {
		return inClusterConfig, nil
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	return clientcmd.BuildConfigFromFlags("", pathOptions.GetDefaultFilename())
}
