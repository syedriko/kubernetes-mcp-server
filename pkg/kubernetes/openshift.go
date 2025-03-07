package kubernetes

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (k *Kubernetes) IsOpenShift(ctx context.Context) bool {
	if _, err := k.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "project.openshift.io",
		Version:  "v1",
		Resource: "projects",
	}).List(ctx, metav1.ListOptions{Limit: 1}); err == nil {
		return true
	}
	return false
}
