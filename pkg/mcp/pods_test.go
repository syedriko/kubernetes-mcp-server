package mcp

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
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
