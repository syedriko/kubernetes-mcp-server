package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"testing"
)

func TestTools(t *testing.T) {
	expectedNames := []string{"pods_list", "pods_list_in_namespace", "configuration_view"}
	t.Run("Has configuration_view tool", testCase(func(t *testing.T, c *mcpContext) {
		tools, err := c.mcpClient.ListTools(c.ctx, mcp.ListToolsRequest{})
		nameSet := make(map[string]bool)
		for _, tool := range tools.Tools {
			nameSet[tool.Name] = true
		}
		for _, name := range expectedNames {
			if nameSet[name] != true {
				t.Fatalf("tool name mismatch %v", err)
				return
			}
		}
	}))
}
