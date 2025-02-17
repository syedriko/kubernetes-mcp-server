package kubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

func (k *Kubernetes) PodsListInAllNamespaces(ctx context.Context) (string, error) {
	return k.ResourcesList(ctx, &schema.GroupVersionKind{
		Group: "", Version: "v1", Kind: "Pod",
	}, "")
}

func (k *Kubernetes) PodsListInNamespace(ctx context.Context, namespace string) (string, error) {
	return k.ResourcesList(ctx, &schema.GroupVersionKind{
		Group: "", Version: "v1", Kind: "Pod",
	}, namespace)
}

func (k *Kubernetes) PodsGet(ctx context.Context, namespace, name string) (string, error) {
	return k.ResourcesGet(ctx, &schema.GroupVersionKind{
		Group: "", Version: "v1", Kind: "Pod",
	}, namespaceOrDefault(namespace), name)
}

func (k *Kubernetes) PodsLog(ctx context.Context, namespace, name string) (string, error) {
	cs, err := kubernetes.NewForConfig(k.cfg)
	if err != nil {
		return "", err
	}
	tailLines := int64(256)
	req := cs.CoreV1().Pods(namespaceOrDefault(namespace)).GetLogs(name, &v1.PodLogOptions{
		TailLines: &tailLines,
	})
	res := req.Do(ctx)
	if res.Error() != nil {
		return "", res.Error()
	}
	rawData, err := res.Raw()
	if err != nil {
		return "", err
	}
	return string(rawData), nil
}
