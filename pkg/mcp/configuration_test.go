package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestConfigurationView(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		toolResult, err := c.callTool("configuration_view", map[string]interface{}{})
		t.Run("configuration_view returns configuration", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
				return
			}
		})
		var decoded *v1.Config
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
		t.Run("configuration_view has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
				return
			}
		})
		t.Run("configuration_view returns current-context", func(t *testing.T) {
			if decoded.CurrentContext != "fake-context" {
				t.Fatalf("fake-context not found: %v", decoded.CurrentContext)
				return
			}
		})
		t.Run("configuration_view returns context info", func(t *testing.T) {
			if len(decoded.Contexts) != 1 {
				t.Fatalf("invalid context count, expected 1, got %v", len(decoded.Contexts))
				return
			}
			if decoded.Contexts[0].Name != "fake-context" {
				t.Fatalf("fake-context not found: %v", decoded.Contexts)
				return
			}
			if decoded.Contexts[0].Context.Cluster != "fake" {
				t.Fatalf("fake-cluster not found: %v", decoded.Contexts)
				return
			}
			if decoded.Contexts[0].Context.AuthInfo != "fake" {
				t.Fatalf("fake-auth not found: %v", decoded.Contexts)
				return
			}
		})
		t.Run("configuration_view returns cluster info", func(t *testing.T) {
			if len(decoded.Clusters) != 1 {
				t.Fatalf("invalid cluster count, expected 1, got %v", len(decoded.Clusters))
				return
			}
			if decoded.Clusters[0].Name != "fake" {
				t.Fatalf("fake-cluster not found: %v", decoded.Clusters)
				return
			}
			if decoded.Clusters[0].Cluster.Server != "https://example.com" {
				t.Fatalf("fake-server not found: %v", decoded.Clusters)
				return
			}
		})
		t.Run("configuration_view returns auth info", func(t *testing.T) {
			if len(decoded.AuthInfos) != 1 {
				t.Fatalf("invalid auth info count, expected 1, got %v", len(decoded.AuthInfos))
				return
			}
			if decoded.AuthInfos[0].Name != "fake" {
				t.Fatalf("fake-auth not found: %v", decoded.AuthInfos)
				return
			}
		})
	})
}
