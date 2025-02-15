package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestPodsListInAllNamespaces(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		createTestData(c.ctx, c.newKubernetesClient())
		configurationGet := mcp.CallToolRequest{}
		configurationGet.Params.Name = "pods_list"
		configurationGet.Params.Arguments = map[string]interface{}{}
		toolResult, err := c.mcpClient.CallTool(c.ctx, configurationGet)
		t.Run("pods_list returns pods list", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
		})
		var decoded []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(map[string]interface{})["text"].(string)), &decoded)
		t.Run("pods_list has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_list returns 2 items", func(t *testing.T) {
			if len(decoded) != 2 {
				t.Fatalf("invalid pods count, expected 2, got %v", len(decoded))
				return
			}
		})
		t.Run("pods_list returns pod in ns-1", func(t *testing.T) {
			if decoded[0].GetName() != "a-pod-in-ns-1" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-1, got %v", decoded[0].GetName())
				return
			}
			if decoded[0].GetNamespace() != "ns-1" {
				t.Fatalf("invalid pod namespace, expected ns-1, got %v", decoded[0].GetNamespace())
				return
			}
		})
		t.Run("pods_list returns pod in ns-2", func(t *testing.T) {
			if decoded[1].GetName() != "a-pod-in-ns-2" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-2, got %v", decoded[1].GetName())
				return
			}
			if decoded[1].GetNamespace() != "ns-2" {
				t.Fatalf("invalid pod namespace, expected ns-2, got %v", decoded[1].GetNamespace())
				return
			}
		})
		t.Run("pods_list omits managed fields", func(t *testing.T) {
			if decoded[0].GetManagedFields() != nil {
				t.Fatalf("managed fields should be omitted, got %v", decoded[0].GetManagedFields())
				return
			}
		})
	})
}

func createTestData(ctx context.Context, kc *kubernetes.Clientset) {
	_, _ = kc.CoreV1().Namespaces().
		Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns-1"}}, metav1.CreateOptions{})
	_, _ = kc.CoreV1().Namespaces().
		Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns-2"}}, metav1.CreateOptions{})
	_, _ = kc.CoreV1().Pods("ns-1").
		Create(ctx, &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "a-pod-in-ns-1"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx"},
				},
			},
		}, metav1.CreateOptions{})
	_, _ = kc.CoreV1().Pods("ns-2").
		Create(ctx, &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "a-pod-in-ns-2"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx"},
				},
			},
		}, metav1.CreateOptions{})
}
