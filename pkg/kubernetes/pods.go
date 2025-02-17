package kubernetes

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	if namespace == "" {
		if ns, _, nsErr := resolveConfig().Namespace(); nsErr == nil {
			namespace = ns
		}
	}
	return k.ResourcesGet(ctx, &schema.GroupVersionKind{
		Group: "", Version: "v1", Kind: "Pod",
	}, namespace, name)
}
