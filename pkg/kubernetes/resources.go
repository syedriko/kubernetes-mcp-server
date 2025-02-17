package kubernetes

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// TODO: WIP
func (k *Kubernetes) ResourcesList(ctx context.Context, gvk *schema.GroupVersionKind, namespace string) (string, error) {
	client, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return "", err
	}
	gvr, err := k.resourceFor(gvk)
	if err != nil {
		return "", err
	}
	rl, err := client.Resource(*gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	return marshal(rl.Items)
}

func (k *Kubernetes) resourceFor(gvk *schema.GroupVersionKind) (*schema.GroupVersionResource, error) {
	if k.deferredDiscoveryRESTMapper == nil {
		d, err := discovery.NewDiscoveryClientForConfig(k.cfg)
		if err != nil {
			return nil, err
		}
		k.deferredDiscoveryRESTMapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(d))
	}
	m, err := k.deferredDiscoveryRESTMapper.RESTMapping(schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}, gvk.Version)
	if err != nil {
		return nil, err
	}
	return &m.Resource, nil
}
