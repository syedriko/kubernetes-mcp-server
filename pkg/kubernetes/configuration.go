package kubernetes

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
)

func ConfigurationView() (string, error) {
	// TODO: consider in cluster run mode (current approach only shows kubeconfig)
	pathOptions := clientcmd.NewDefaultPathOptions()
	cfg, err := pathOptions.GetStartingConfig()
	if err != nil {
		return "", err
	}
	if err = clientcmdapi.MinifyConfig(cfg); err != nil {
		return "", err
	}
	if err = clientcmdapi.FlattenConfig(cfg); err != nil {
		return "", err
	}
	convertedObj, err := latest.Scheme.ConvertToVersion(cfg, latest.ExternalVersion)
	if err != nil {
		return "", err
	}
	return marshal(convertedObj)
}
