package mcp

import (
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"testing"
)

func TestResourcesCreateOrUpdate(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		t.Run("resources_create_or_update with nil resource returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("resources_create_or_update", map[string]interface{}{})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to create or update resources, missing argument resource" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		t.Run("resources_create_or_update with empty resource returns error", func(t *testing.T) {
			toolResult, _ := c.callTool("resources_create_or_update", map[string]interface{}{"resource": ""})
			if toolResult.IsError != true {
				t.Fatalf("call tool should fail")
				return
			}
			if toolResult.Content[0].(map[string]interface{})["text"].(string) != "failed to create or update resources, missing argument resource" {
				t.Fatalf("invalid error message, got %v", toolResult.Content[0].(map[string]interface{})["text"].(string))
				return
			}
		})
		client := c.newKubernetesClient()
		configMapYaml := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: a-cm-created-or-updated\n  namespace: default\n"
		resourcesCreateOrUpdateCm1, err := c.callTool("resources_create_or_update", map[string]interface{}{"resource": configMapYaml})
		t.Run("resources_create_or_update with valid namespaced yaml resource returns success", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if resourcesCreateOrUpdateCm1.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		t.Run("resources_create_or_update with valid namespaced yaml resource creates ConfigMap", func(t *testing.T) {
			cm, _ := client.CoreV1().ConfigMaps("default").Get(c.ctx, "a-cm-created-or-updated", metav1.GetOptions{})
			if cm == nil {
				t.Fatalf("ConfigMap not found")
				return
			}
		})
		configMapJson := "{\"apiVersion\": \"v1\", \"kind\": \"ConfigMap\", \"metadata\": {\"name\": \"a-cm-created-or-updated-2\", \"namespace\": \"default\"}}"
		resourcesCreateOrUpdateCm2, err := c.callTool("resources_create_or_update", map[string]interface{}{"resource": configMapJson})
		t.Run("resources_create_or_update with valid namespaced json resource returns success", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if resourcesCreateOrUpdateCm2.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		t.Run("resources_create_or_update with valid namespaced json resource creates config map", func(t *testing.T) {
			cm, _ := client.CoreV1().ConfigMaps("default").Get(c.ctx, "a-cm-created-or-updated-2", metav1.GetOptions{})
			if cm == nil {
				t.Fatalf("ConfigMap not found")
				return
			}
		})
		customResourceDefinitionJson := `
          {
            "apiVersion": "apiextensions.k8s.io/v1",
            "kind": "CustomResourceDefinition",
            "metadata": {"name": "customs.example.com"},
            "spec": {
              "group": "example.com",
              "versions": [{
                "name": "v1","served": true,"storage": true,
                "schema": {"openAPIV3Schema": {"type": "object"}}
              }],
              "scope": "Namespaced",
              "names": {"plural": "customs","singular": "custom","kind": "Custom"}
            }
          }`
		resourcesCreateOrUpdateCrd, err := c.callTool("resources_create_or_update", map[string]interface{}{"resource": customResourceDefinitionJson})
		t.Run("resources_create_or_update with valid cluster-scoped json resource returns success", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if resourcesCreateOrUpdateCrd.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		t.Run("resources_create_or_update with valid cluster-scoped json resource creates custom resource definition", func(t *testing.T) {
			apiExtensionsV1Client := v1.NewForConfigOrDie(envTestRestConfig)
			_, err = apiExtensionsV1Client.CustomResourceDefinitions().Get(c.ctx, "customs.example.com", metav1.GetOptions{})
			if err != nil {
				t.Fatalf("custom resource definition not found")
				return
			}
		})
		customJson := "{\"apiVersion\": \"example.com/v1\", \"kind\": \"Custom\", \"metadata\": {\"name\": \"a-custom-resource\"}}"
		resourcesCreateOrUpdateCustom, err := c.callTool("resources_create_or_update", map[string]interface{}{"resource": customJson})
		t.Run("resources_create_or_update with valid namespaced json resource returns success", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if resourcesCreateOrUpdateCustom.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		t.Run("resources_create_or_update with valid namespaced json resource creates custom resource", func(t *testing.T) {
			dynamicClient := dynamic.NewForConfigOrDie(envTestRestConfig)
			_, err = dynamicClient.
				Resource(schema.GroupVersionResource{Group: "example.com", Version: "v1", Resource: "customs"}).
				Namespace("default").
				Get(c.ctx, "a-custom-resource", metav1.GetOptions{})
			if err != nil {
				t.Fatalf("custom resource not found")
				return
			}
		})
		customJsonUpdated := "{\"apiVersion\": \"example.com/v1\", \"kind\": \"Custom\", \"metadata\": {\"name\": \"a-custom-resource\",\"annotations\": {\"updated\": \"true\"}}}"
		resourcesCreateOrUpdateCustomUpdated, err := c.callTool("resources_create_or_update", map[string]interface{}{"resource": customJsonUpdated})
		t.Run("resources_create_or_update with valid namespaced json resource updates custom resource", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
			if resourcesCreateOrUpdateCustomUpdated.IsError {
				t.Fatalf("call tool failed")
				return
			}
		})
		t.Run("resources_create_or_update with valid namespaced json resource updates custom resource", func(t *testing.T) {
			dynamicClient := dynamic.NewForConfigOrDie(envTestRestConfig)
			customResource, _ := dynamicClient.
				Resource(schema.GroupVersionResource{Group: "example.com", Version: "v1", Resource: "customs"}).
				Namespace("default").
				Get(c.ctx, "a-custom-resource", metav1.GetOptions{})
			if customResource == nil {
				t.Fatalf("custom resource not found")
				return
			}
			annotations := customResource.GetAnnotations()
			if annotations == nil || annotations["updated"] != "true" {
				t.Fatalf("custom resource not updated")
				return
			}
		})
	})
}
