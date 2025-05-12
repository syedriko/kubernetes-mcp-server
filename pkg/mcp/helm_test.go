package mcp

import (
	"encoding/base64"
	"github.com/mark3labs/mcp-go/mcp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/yaml"
	"strings"
	"testing"
)

func TestHelmInstall(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		_, file, _, _ := runtime.Caller(0)
		chartPath := filepath.Join(filepath.Dir(file), "testdata", "helm-chart-no-op")
		toolResult, err := c.callTool("helm_install", map[string]interface{}{
			"chart": chartPath,
		})
		t.Run("helm_install with local chart and no release name, returns installed chart", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			var decoded []map[string]interface{}
			err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
			if !strings.HasPrefix(decoded[0]["name"].(string), "helm-chart-no-op-") {
				t.Fatalf("invalid helm install name, expected no-op-*, got %v", decoded[0]["name"])
			}
			if decoded[0]["namespace"] != "default" {
				t.Fatalf("invalid helm install namespace, expected default, got %v", decoded[0]["namespace"])
			}
			if decoded[0]["chart"] != "no-op" {
				t.Fatalf("invalid helm install name, expected release name, got empty")
			}
			if decoded[0]["chartVersion"] != "1.33.7" {
				t.Fatalf("invalid helm install version, expected 1.33.7, got empty")
			}
			if decoded[0]["status"] != "deployed" {
				t.Fatalf("invalid helm install status, expected deployed, got %v", decoded[0]["status"])
			}
			if decoded[0]["revision"] != float64(1) {
				t.Fatalf("invalid helm install revision, expected 1, got %v", decoded[0]["revision"])
			}
		})
	})
}

func TestHelmList(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		c.withEnvTest()
		kc := c.newKubernetesClient()
		secrets, err := kc.CoreV1().Secrets("default").List(c.ctx, metav1.ListOptions{})
		for _, secret := range secrets.Items {
			if strings.HasPrefix(secret.Name, "sh.helm.release.v1.") {
				_ = kc.CoreV1().Secrets("default").Delete(c.ctx, secret.Name, metav1.DeleteOptions{})
			}
		}
		_ = kc.CoreV1().Secrets("default").Delete(c.ctx, "release-to-list", metav1.DeleteOptions{})
		toolResult, err := c.callTool("helm_list", map[string]interface{}{})
		t.Run("helm_list with no releases, returns not found", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			if toolResult.Content[0].(mcp.TextContent).Text != "No Helm releases found" {
				t.Fatalf("unexpected result %v", toolResult.Content[0].(mcp.TextContent).Text)
			}
		})
		_, _ = kc.CoreV1().Secrets("default").Create(c.ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "release-to-list",
				Labels: map[string]string{"owner": "helm"},
			},
			Data: map[string][]byte{
				"release": []byte(base64.StdEncoding.EncodeToString([]byte("{" +
					"\"name\":\"release-to-list\"," +
					"\"info\":{\"status\":\"deployed\"}" +
					"}"))),
			},
		}, metav1.CreateOptions{})
		toolResult, err = c.callTool("helm_list", map[string]interface{}{})
		t.Run("helm_list with deployed release, returns release", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			var decoded []map[string]interface{}
			err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
			if len(decoded) != 1 {
				t.Fatalf("invalid helm list count, expected 1, got %v", len(decoded))
			}
			if decoded[0]["name"] != "release-to-list" {
				t.Fatalf("invalid helm list name, expected release-to-list, got %v", decoded[0]["name"])
			}
			if decoded[0]["status"] != "deployed" {
				t.Fatalf("invalid helm list status, expected deployed, got %v", decoded[0]["status"])
			}
		})
		toolResult, err = c.callTool("helm_list", map[string]interface{}{"namespace": "ns-1"})
		t.Run("helm_list with deployed release in other namespaces, returns not found", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			if toolResult.Content[0].(mcp.TextContent).Text != "No Helm releases found" {
				t.Fatalf("unexpected result %v", toolResult.Content[0].(mcp.TextContent).Text)
			}
		})
		toolResult, err = c.callTool("helm_list", map[string]interface{}{"namespace": "ns-1", "all_namespaces": true})
		t.Run("helm_list with deployed release in all namespaces, returns release", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			var decoded []map[string]interface{}
			err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
			if len(decoded) != 1 {
				t.Fatalf("invalid helm list count, expected 1, got %v", len(decoded))
			}
			if decoded[0]["name"] != "release-to-list" {
				t.Fatalf("invalid helm list name, expected release-to-list, got %v", decoded[0]["name"])
			}
			if decoded[0]["status"] != "deployed" {
				t.Fatalf("invalid helm list status, expected deployed, got %v", decoded[0]["status"])
			}
		})
	})
}
