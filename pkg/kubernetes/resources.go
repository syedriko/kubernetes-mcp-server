package kubernetes

import (
	"context"
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"regexp"
	"strings"
)

const (
	AppKubernetesComponent = "app.kubernetes.io/component"
	AppKubernetesManagedBy = "app.kubernetes.io/managed-by"
	AppKubernetesName      = "app.kubernetes.io/name"
	AppKubernetesPartOf    = "app.kubernetes.io/part-of"
)

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

func (k *Kubernetes) ResourcesGet(ctx context.Context, gvk *schema.GroupVersionKind, namespace, name string) (string, error) {
	client, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return "", err
	}
	gvr, err := k.resourceFor(gvk)
	if err != nil {
		return "", err
	}
	// If it's a namespaced resource and namespace wasn't provided, try to use the default configured one
	if namespaced, nsErr := k.isNamespaced(gvk); nsErr == nil && namespaced {
		namespace = namespaceOrDefault(namespace)
	}
	rg, err := client.Resource(*gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return marshal(rg)
}

func (k *Kubernetes) ResourcesCreateOrUpdate(ctx context.Context, resource string) (string, error) {
	separator := regexp.MustCompile(`\r?\n---\r?\n`)
	resources := separator.Split(resource, -1)
	var parsedResources []*unstructured.Unstructured
	for _, r := range resources {
		var obj unstructured.Unstructured
		if err := yaml.NewYAMLToJSONDecoder(strings.NewReader(r)).Decode(&obj); err != nil {
			return "", err
		}
		parsedResources = append(parsedResources, &obj)
	}
	return k.resourcesCreateOrUpdate(ctx, parsedResources)
}

func (k *Kubernetes) ResourcesDelete(ctx context.Context, gvk *schema.GroupVersionKind, namespace, name string) error {
	client, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return err
	}
	gvr, err := k.resourceFor(gvk)
	if err != nil {
		return err
	}
	// If it's a namespaced resource and namespace wasn't provided, try to use the default configured one
	if namespaced, nsErr := k.isNamespaced(gvk); nsErr == nil && namespaced {
		namespace = namespaceOrDefault(namespace)
	}
	return client.Resource(*gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (k *Kubernetes) resourcesCreateOrUpdate(ctx context.Context, resources []*unstructured.Unstructured) (string, error) {
	client, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return "", err
	}
	for i, obj := range resources {
		gvk := obj.GroupVersionKind()
		gvr, rErr := k.resourceFor(&gvk)
		if rErr != nil {
			return "", rErr
		}
		namespace := obj.GetNamespace()
		// If it's a namespaced resource and namespace wasn't provided, try to use the default configured one
		if namespaced, nsErr := k.isNamespaced(&gvk); nsErr == nil && namespaced {
			namespace = namespaceOrDefault(namespace)
		}
		resources[i], rErr = client.Resource(*gvr).Namespace(namespace).Apply(ctx, obj.GetName(), obj, metav1.ApplyOptions{
			FieldManager: version.BinaryName,
		})
		if rErr != nil {
			return "", rErr
		}
	}
	return marshal(resources)
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

func (k *Kubernetes) isNamespaced(gvk *schema.GroupVersionKind) (bool, error) {
	d, err := discovery.NewDiscoveryClientForConfig(k.cfg)
	if err != nil {
		return false, err
	}
	apiResourceList, err := d.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return false, err
	}
	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == gvk.Kind {
			return apiResource.Namespaced, nil
		}
	}
	return false, nil
}

func (k *Kubernetes) supportsGroupVersion(groupVersion string) bool {
	d, err := discovery.NewDiscoveryClientForConfig(k.cfg)
	if err != nil {
		return false
	}
	_, err = d.ServerResourcesForGroupVersion(groupVersion)
	if err == nil {
		return true
	}
	return false
}
