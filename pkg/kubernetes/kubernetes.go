package kubernetes

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
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

func defaultPrintFlags() *genericclioptions.PrintFlags {
	return genericclioptions.NewPrintFlags("").
		WithTypeSetter(scheme.Scheme).
		WithDefaultOutput("yaml")
}

func resolveClientConfig() (*rest.Config, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err == nil && inClusterConfig != nil {
		return inClusterConfig, nil
	}
	pathOptions := clientcmd.NewDefaultPathOptions()
	return clientcmd.BuildConfigFromFlags("", pathOptions.GetDefaultFilename())
}
