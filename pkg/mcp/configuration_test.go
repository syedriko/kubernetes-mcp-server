package mcp

import (
	"github.com/manusa/kubernetes-mcp-server/pkg/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
	"k8s.io/client-go/rest"
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
			}
		})
		var decoded *v1.Config
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
		t.Run("configuration_view has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
		})
		t.Run("configuration_view returns current-context", func(t *testing.T) {
			if decoded.CurrentContext != "fake-context" {
				t.Fatalf("fake-context not found: %v", decoded.CurrentContext)
			}
		})
		t.Run("configuration_view returns context info", func(t *testing.T) {
			if len(decoded.Contexts) != 1 {
				t.Fatalf("invalid context count, expected 1, got %v", len(decoded.Contexts))
			}
			if decoded.Contexts[0].Name != "fake-context" {
				t.Fatalf("fake-context not found: %v", decoded.Contexts)
			}
			if decoded.Contexts[0].Context.Cluster != "fake" {
				t.Fatalf("fake-cluster not found: %v", decoded.Contexts)
			}
			if decoded.Contexts[0].Context.AuthInfo != "fake" {
				t.Fatalf("fake-auth not found: %v", decoded.Contexts)
			}
		})
		t.Run("configuration_view returns cluster info", func(t *testing.T) {
			if len(decoded.Clusters) != 1 {
				t.Fatalf("invalid cluster count, expected 1, got %v", len(decoded.Clusters))
			}
			if decoded.Clusters[0].Name != "fake" {
				t.Fatalf("fake-cluster not found: %v", decoded.Clusters)
			}
			if decoded.Clusters[0].Cluster.Server != "https://example.com" {
				t.Fatalf("fake-server not found: %v", decoded.Clusters)
			}
		})
		t.Run("configuration_view returns auth info", func(t *testing.T) {
			if len(decoded.AuthInfos) != 1 {
				t.Fatalf("invalid auth info count, expected 1, got %v", len(decoded.AuthInfos))
			}
			if decoded.AuthInfos[0].Name != "fake" {
				t.Fatalf("fake-auth not found: %v", decoded.AuthInfos)
			}
		})
	})
}

func TestConfigurationViewInCluster(t *testing.T) {
	kubernetes.InClusterConfig = func() (*rest.Config, error) {
		return &rest.Config{
			Host:        "https://kubernetes.default.svc",
			BearerToken: "fake-token",
		}, nil
	}
	defer func() {
		kubernetes.InClusterConfig = rest.InClusterConfig
	}()
	testCase(t, func(c *mcpContext) {
		toolResult, err := c.callTool("configuration_view", map[string]interface{}{})
		t.Run("configuration_view returns configuration", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
		})
		var decoded *v1.Config
		err = yaml.Unmarshal([]byte(toolResult.Content[0].(mcp.TextContent).Text), &decoded)
		t.Run("configuration_view has yaml content", func(t *testing.T) {
			if err != nil {
				t.Fatalf("invalid tool result content %v", err)
			}
		})
		t.Run("configuration_view returns current-context", func(t *testing.T) {
			if decoded.CurrentContext != "context" {
				t.Fatalf("context not found: %v", decoded.CurrentContext)
			}
		})
		t.Run("configuration_view returns context info", func(t *testing.T) {
			if len(decoded.Contexts) != 1 {
				t.Fatalf("invalid context count, expected 1, got %v", len(decoded.Contexts))
			}
			if decoded.Contexts[0].Name != "context" {
				t.Fatalf("context not found: %v", decoded.Contexts)
			}
			if decoded.Contexts[0].Context.Cluster != "cluster" {
				t.Fatalf("cluster not found: %v", decoded.Contexts)
			}
			if decoded.Contexts[0].Context.AuthInfo != "user" {
				t.Fatalf("user not found: %v", decoded.Contexts)
			}
		})
		t.Run("configuration_view returns cluster info", func(t *testing.T) {
			if len(decoded.Clusters) != 1 {
				t.Fatalf("invalid cluster count, expected 1, got %v", len(decoded.Clusters))
			}
			if decoded.Clusters[0].Name != "cluster" {
				t.Fatalf("cluster not found: %v", decoded.Clusters)
			}
			if decoded.Clusters[0].Cluster.Server != "https://kubernetes.default.svc" {
				t.Fatalf("server not found: %v", decoded.Clusters)
			}
		})
		t.Run("configuration_view returns auth info", func(t *testing.T) {
			if len(decoded.AuthInfos) != 1 {
				t.Fatalf("invalid auth info count, expected 1, got %v", len(decoded.AuthInfos))
			}
			if decoded.AuthInfos[0].Name != "user" {
				t.Fatalf("user not found: %v", decoded.AuthInfos)
			}
		})
	})
}
