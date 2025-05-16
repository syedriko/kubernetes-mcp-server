package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/yaml"
	"slices"
	"testing"
)

func TestNamespacesList(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		toolResult, err := c.callTool("namespaces_list", map[string]interface{}{})
		t.Run("namespaces_list returns namespace list", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
		})
		var decoded []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
		t.Run("namespaces_list has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
		})
		t.Run("namespaces_list returns at least 3 items", func(t *testing.T) {
			if len(decoded) < 3 {
				t.Errorf("invalid namespace count, expected at least 3, got %v", len(decoded))
			}
			for _, expectedNamespace := range []string{"default", "ns-1", "ns-2"} {
				idx := slices.IndexFunc(decoded, func(ns unstructured.Unstructured) bool {
					return ns.GetName() == expectedNamespace
				})
				if idx == -1 {
					t.Errorf("namespace %s not found in the list", expectedNamespace)
				}
			}
		})
	})
}

func TestProjectsListInOpenShift(t *testing.T) {
	testCaseWithContext(t, &mcpContext{before: inOpenShift, after: inOpenShiftClear}, func(c *mcpContext) {
		dynamicClient := dynamic.NewForConfigOrDie(envTestRestConfig)
		_, _ = dynamicClient.Resource(schema.GroupVersionResource{Group: "project.openshift.io", Version: "v1", Resource: "projects"}).
			Create(c.ctx, &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": "project.openshift.io/v1",
				"kind":       "Project",
				"metadata": map[string]interface{}{
					"name": "an-openshift-project",
				},
			}}, metav1.CreateOptions{})
		toolResult, err := c.callTool("projects_list", map[string]interface{}{})
		t.Run("projects_list returns project list", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
		})
		var decoded []unstructured.Unstructured
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
		t.Run("projects_list has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
		})
		t.Run("projects_list returns at least 1 items", func(t *testing.T) {
			if len(decoded) < 1 {
				t.Errorf("invalid project count, expected at least 1, got %v", len(decoded))
			}
			idx := slices.IndexFunc(decoded, func(ns unstructured.Unstructured) bool {
				return ns.GetName() == "an-openshift-project"
			})
			if idx == -1 {
				t.Errorf("namespace %s not found in the list", "an-openshift-project")
			}
		})
	})
}
