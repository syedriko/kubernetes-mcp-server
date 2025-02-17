package mcp

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
	"strings"
	"testing"
)

func TestPodsListInAllNamespaces(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		toolResult, err := c.callTool("pods_list", map[string]interface{}{})
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
		t.Run("pods_list returns 3 items", func(t *testing.T) {
			if len(decoded) != 3 {
				t.Fatalf("invalid pods count, expected 3, got %v", len(decoded))
				return
			}
		})
		t.Run("pods_list returns pod in ns-1", func(t *testing.T) {
			if decoded[1].GetName() != "a-pod-in-ns-1" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-1, got %v", decoded[1].GetName())
				return
			}
			if decoded[1].GetNamespace() != "ns-1" {
				t.Fatalf("invalid pod namespace, expected ns-1, got %v", decoded[1].GetNamespace())
				return
			}
		})
		t.Run("pods_list returns pod in ns-2", func(t *testing.T) {
			if decoded[2].GetName() != "a-pod-in-ns-2" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-2, got %v", decoded[2].GetName())
				return
			}
			if decoded[2].GetNamespace() != "ns-2" {
				t.Fatalf("invalid pod namespace, expected ns-2, got %v", decoded[2].GetNamespace())
				return
			}
		})
		t.Run("pods_list omits managed fields", func(t *testing.T) {
			if decoded[1].GetManagedFields() != nil {
				t.Fatalf("managed fields should be omitted, got %v", decoded[0].GetManagedFields())
				return
			}
		})
	})
}

func TestPodsListInNamespace(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		t.Run("pods_list_in_namespace with nil namespace returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("pods_list_in_namespace", map[string]interface{}{})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to list pods in namespace, missing argument namespace" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		toolResult, err := c.callTool("pods_list_in_namespace", map[string]interface{}{
			"namespace": "ns-1",
		})
		t.Run("pods_list_in_namespace returns pods list", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		var decoded []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(map[string]interface{})["text"].(string)), &decoded)
		t.Run("pods_list_in_namespace has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_list_in_namespace returns 1 items", func(t *testing.T) {
			if len(decoded) != 1 {
				t.Fatalf("invalid pods count, expected 1, got %v", len(decoded))
				return
			}
		})
		t.Run("pods_list_in_namespace returns pod in ns-1", func(t *testing.T) {
			if decoded[0].GetName() != "a-pod-in-ns-1" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-1, got %v", decoded[0].GetName())
				return
			}
			if decoded[0].GetNamespace() != "ns-1" {
				t.Fatalf("invalid pod namespace, expected ns-1, got %v", decoded[0].GetNamespace())
				return
			}
		})
		t.Run("pods_list_in_namespace omits managed fields", func(t *testing.T) {
			if decoded[0].GetManagedFields() != nil {
				t.Fatalf("managed fields should be omitted, got %v", decoded[0].GetManagedFields())
				return
			}
		})
	})
}

func TestPodsGet(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		t.Run("pods_get with nil name returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("pods_get", map[string]interface{}{})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to get pod, missing argument name" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		podsGetNilNamespace, err := c.callTool("pods_get", map[string]interface{}{
			"name": "a-pod-in-default",
		})
		t.Run("pods_get with name and nil namespace returns pod", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if podsGetNilNamespace.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		var decodedNilNamespace unstructured.Unstructured
		err = yaml.Unmarshal([]byte(podsGetNilNamespace.Content[0].(map[string]interface{})["text"].(string)), &decodedNilNamespace)
		t.Run("pods_get with name and nil namespace has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_get with name and nil namespace returns pod in default", func(t *testing.T) {
			if decodedNilNamespace.GetName() != "a-pod-in-default" {
				t.Fatalf("invalid pod name, expected a-pod-in-default, got %v", decodedNilNamespace.GetName())
				return
			}
			if decodedNilNamespace.GetNamespace() != "default" {
				t.Fatalf("invalid pod namespace, expected default, got %v", decodedNilNamespace.GetNamespace())
				return
			}
		})
		podsGetInNamespace, err := c.callTool("pods_get", map[string]interface{}{
			"namespace": "ns-1",
			"name":      "a-pod-in-ns-1",
		})
		t.Run("pods_get with name and namespace returns pod", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if podsGetInNamespace.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		var decodedInNamespace unstructured.Unstructured
		err = yaml.Unmarshal([]byte(podsGetInNamespace.Content[0].(map[string]interface{})["text"].(string)), &decodedInNamespace)
		t.Run("pods_get with name and namespace has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_get with name and namespace returns pod in ns-1", func(t *testing.T) {
			if decodedInNamespace.GetName() != "a-pod-in-ns-1" {
				t.Fatalf("invalid pod name, expected a-pod-in-ns-1, got %v", decodedInNamespace.GetName())
				return
			}
			if decodedInNamespace.GetNamespace() != "ns-1" {
				t.Fatalf("invalid pod namespace, ns-1 ns-1, got %v", decodedInNamespace.GetNamespace())
				return
			}
		})
	})
}

func TestPodsLog(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		t.Run("pods_log with nil name returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("pods_log", map[string]interface{}{})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to get pod log, missing argument name" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		t.Run("pods_log with not found name returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("pods_log", map[string]interface{}{"name": "not-found"})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to get pod not-found log in namespace : pods \"not-found\" not found" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		podsLogNilNamespace, err := c.callTool("pods_log", map[string]interface{}{
			"name": "a-pod-in-default",
		})
		t.Run("pods_log with name and nil namespace returns pod log", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if podsLogNilNamespace.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		podsLogInNamespace, err := c.callTool("pods_log", map[string]interface{}{
			"namespace": "ns-1",
			"name":      "a-pod-in-ns-1",
		})
		t.Run("pods_log with name and namespace returns pod log", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if podsLogInNamespace.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
	})
}

func TestPodsRun(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		t.Run("pods_run with nil image returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("pods_run", map[string]interface{}{})
			if toolResult.IsError != true {
				t.Errorf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to run pod, missing argument image" {
				t.Errorf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		podsRunNilNamespace, err := c.callTool("pods_run", map[string]interface{}{"image": "nginx"})
		t.Run("pods_run with image and nil namespace runs pod", func(t *testing.T) {
			if err != nil {
				t.Errorf("call tool failed %v", err)
				return
			}
			if podsRunNilNamespace.IsError {
				t.Errorf("call tool failed")
				return
			}
		})
		var decodedNilNamespace []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(podsRunNilNamespace.Content[0].(map[string]interface{})["text"].(string)), &decodedNilNamespace)
		t.Run("pods_run with image and nil namespace has yaml content", func(t *testing.T) {
			if err != nil {
				t.Errorf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_run with image and nil namespace returns 1 item (Pod)", func(t *testing.T) {
			if len(decodedNilNamespace) != 1 {
				t.Errorf("invalid pods count, expected 1, got %v", len(decodedNilNamespace))
				return
			}
			if decodedNilNamespace[0].GetKind() != "Pod" {
				t.Errorf("invalid pod kind, expected Pod, got %v", decodedNilNamespace[0].GetKind())
				return
			}
		})
		t.Run("pods_run with image and nil namespace returns pod in default", func(t *testing.T) {
			if decodedNilNamespace[0].GetNamespace() != "default" {
				t.Errorf("invalid pod namespace, expected default, got %v", decodedNilNamespace[0].GetNamespace())
				return
			}
		})
		t.Run("pods_run with image and nil namespace returns pod with random name", func(t *testing.T) {
			if !strings.HasPrefix(decodedNilNamespace[0].GetName(), "kubernetes-mcp-server-run-") {
				t.Errorf("invalid pod name, expected random, got %v", decodedNilNamespace[0].GetName())
				return
			}
		})
		t.Run("pods_run with image and nil namespace returns pod with labels", func(t *testing.T) {
			labels := decodedNilNamespace[0].Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})
			if labels["app.kubernetes.io/name"] == "" {
				t.Errorf("invalid labels, expected app.kubernetes.io/name, got %v", labels)
				return
			}
			if labels["app.kubernetes.io/component"] == "" {
				t.Errorf("invalid labels, expected app.kubernetes.io/component, got %v", labels)
				return
			}
			if labels["app.kubernetes.io/managed-by"] != "kubernetes-mcp-server" {
				t.Errorf("invalid labels, expected app.kubernetes.io/managed-by, got %v", labels)
				return
			}
			if labels["app.kubernetes.io/part-of"] != "kubernetes-mcp-server-run-sandbox" {
				t.Errorf("invalid labels, expected app.kubernetes.io/part-of, got %v", labels)
				return
			}
		})
		t.Run("pods_run with image and nil namespace returns pod with nginx container", func(t *testing.T) {
			containers := decodedNilNamespace[0].Object["spec"].(map[string]interface{})["containers"].([]interface{})
			if containers[0].(map[string]interface{})["image"] != "nginx" {
				t.Errorf("invalid container name, expected nginx, got %v", containers[0].(map[string]interface{})["image"])
				return
			}
		})

		podsRunNamespaceAndPort, err := c.callTool("pods_run", map[string]interface{}{"image": "nginx", "port": 80})
		t.Run("pods_run with image, namespace, and port runs pod", func(t *testing.T) {
			if err != nil {
				t.Errorf("call tool failed %v", err)
				return
			}
			if podsRunNamespaceAndPort.IsError {
				t.Errorf("call tool failed")
				return
			}
		})
		var decodedNamespaceAndPort []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(podsRunNamespaceAndPort.Content[0].(map[string]interface{})["text"].(string)), &decodedNamespaceAndPort)
		t.Run("pods_run with image, namespace, and port has yaml content", func(t *testing.T) {
			if err != nil {
				t.Errorf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("pods_run with image, namespace, and port returns 2 items (Pod + Service)", func(t *testing.T) {
			if len(decodedNamespaceAndPort) != 2 {
				t.Errorf("invalid pods count, expected 2, got %v", len(decodedNamespaceAndPort))
				return
			}
			if decodedNamespaceAndPort[0].GetKind() != "Pod" {
				t.Errorf("invalid pod kind, expected Pod, got %v", decodedNamespaceAndPort[0].GetKind())
				return
			}
			if decodedNamespaceAndPort[1].GetKind() != "Service" {
				t.Errorf("invalid service kind, expected Service, got %v", decodedNamespaceAndPort[1].GetKind())
				return
			}
		})
		t.Run("pods_run with image, namespace, and port returns pod with port", func(t *testing.T) {
			containers := decodedNamespaceAndPort[0].Object["spec"].(map[string]interface{})["containers"].([]interface{})
			ports := containers[0].(map[string]interface{})["ports"].([]interface{})
			if ports[0].(map[string]interface{})["containerPort"] != int64(80) {
				t.Errorf("invalid container port, expected 80, got %v", ports[0].(map[string]interface{})["containerPort"])
				return
			}
		})
		t.Run("pods_run with image, namespace, and port returns service with port and selector", func(t *testing.T) {
			ports := decodedNamespaceAndPort[1].Object["spec"].(map[string]interface{})["ports"].([]interface{})
			if ports[0].(map[string]interface{})["port"] != int64(80) {
				t.Errorf("invalid service port, expected 80, got %v", ports[0].(map[string]interface{})["port"])
				return
			}
			if ports[0].(map[string]interface{})["targetPort"] != int64(80) {
				t.Errorf("invalid service target port, expected 80, got %v", ports[0].(map[string]interface{})["targetPort"])
				return
			}
			selector := decodedNamespaceAndPort[1].Object["spec"].(map[string]interface{})["selector"].(map[string]interface{})
			if selector["app.kubernetes.io/name"] == "" {
				t.Errorf("invalid service selector, expected app.kubernetes.io/name, got %v", selector)
				return
			}
			if selector["app.kubernetes.io/managed-by"] != "kubernetes-mcp-server" {
				t.Errorf("invalid service selector, expected app.kubernetes.io/managed-by, got %v", selector)
				return
			}
			if selector["app.kubernetes.io/part-of"] != "kubernetes-mcp-server-run-sandbox" {
				t.Errorf("invalid service selector, expected app.kubernetes.io/part-of, got %v", selector)
				return
			}
		})
	})
}
