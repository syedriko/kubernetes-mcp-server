package kubernetes

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
)

func ConfigurationView(minify bool) (string, error) {
	var cfg clientcmdapi.Config
	var err error
	inClusterConfig, err := InClusterConfig()
	if err == nil && inClusterConfig != nil {
		cfg = *clientcmdapi.NewConfig()
		cfg.Clusters["cluster"] = &clientcmdapi.Cluster{
			Server:                inClusterConfig.Host,
			InsecureSkipTLSVerify: inClusterConfig.Insecure,
		}
		cfg.AuthInfos["user"] = &clientcmdapi.AuthInfo{
			Token: inClusterConfig.BearerToken,
		}
		cfg.Contexts["context"] = &clientcmdapi.Context{
			Cluster:  "cluster",
			AuthInfo: "user",
		}
		cfg.CurrentContext = "context"
	} else if cfg, err = resolveConfig().RawConfig(); err != nil {
		return "", err
	}
	if minify {
		if err = clientcmdapi.MinifyConfig(&cfg); err != nil {
			return "", err
		}
	}
	if err = clientcmdapi.FlattenConfig(&cfg); err != nil {
		return "", err
	}
	convertedObj, err := latest.Scheme.ConvertToVersion(&cfg, latest.ExternalVersion)
	if err != nil {
		return "", err
	}
	return marshal(convertedObj)
}
