package kubernetes

import (
	"context"
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
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

func (k *Kubernetes) PodsDelete(ctx context.Context, namespace, name string) (string, error) {
	// TODO
	return "", nil
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

func (k *Kubernetes) PodsRun(ctx context.Context, namespace, name, image string, port int32) (string, error) {
	if name == "" {
		name = version.BinaryName + "-run-" + rand.String(5)
	}
	labels := map[string]string{
		AppKubernetesName:      name,
		AppKubernetesComponent: name,
		AppKubernetesManagedBy: version.BinaryName,
		AppKubernetesPartOf:    version.BinaryName + "-run-sandbox",
	}
	// NewPod
	var resources []any
	pod := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceOrDefault(namespace), Labels: labels},
		Spec: v1.PodSpec{Containers: []v1.Container{{
			Name:            name,
			Image:           image,
			ImagePullPolicy: v1.PullAlways,
		}}},
	}
	resources = append(resources, pod)
	if port > 0 {
		pod.Spec.Containers[0].Ports = []v1.ContainerPort{{ContainerPort: port}}
		resources = append(resources, &v1.Service{
			TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespaceOrDefault(namespace), Labels: labels},
			Spec: v1.ServiceSpec{
				Selector: labels,
				Type:     v1.ServiceTypeClusterIP,
				Ports:    []v1.ServicePort{{Port: port, TargetPort: intstr.FromInt32(port)}},
			},
		})
	}
	if port > 0 && k.supportsGroupVersion("route.openshift.io/v1") {
		resources = append(resources, &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "route.openshift.io/v1",
				"kind":       "Route",
				"metadata": map[string]interface{}{
					"name":      name,
					"namespace": namespaceOrDefault(namespace),
					"labels":    labels,
				},
				"spec": map[string]interface{}{
					"to": map[string]interface{}{
						"kind":   "Service",
						"name":   name,
						"weight": 100,
					},
					"port": map[string]interface{}{
						"targetPort": intstr.FromInt32(port),
					},
					"tls": map[string]interface{}{
						"termination":                   "edge",
						"insecureEdgeTerminationPolicy": "Redirect",
					},
				},
			},
		})

	}

	// Convert the objects to Unstructured and reuse resourcesCreateOrUpdate functionality
	converter := runtime.DefaultUnstructuredConverter
	var toCreate []*unstructured.Unstructured
	for _, obj := range resources {
		m, err := converter.ToUnstructured(obj)
		if err != nil {
			return "", err
		}
		u := &unstructured.Unstructured{}
		if err = converter.FromUnstructured(m, u); err != nil {
			return "", err
		}
		toCreate = append(toCreate, u)
	}
	return k.resourcesCreateOrUpdate(ctx, toCreate)
}
