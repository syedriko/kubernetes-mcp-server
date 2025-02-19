package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"testing"
)

func TestTools(t *testing.T) {
	expectedNames := []string{
		"configuration_view",
		"pods_list",
		"pods_list_in_namespace",
		"pods_get",
		"pods_delete",
		"pods_log",
		"pods_run",
		"resources_list",
		"resources_get",
		"resources_create_or_update",
		"resources_delete",
	}
	testCase(t, func(c *mcpContext) {
		tools, err := c.mcpClient.ListTools(c.ctx, mcp.ListToolsRequest{})
		t.Run("ListTools returns tools", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call ListTools failed %v", err)
				return
			}
		})
		nameSet := make(map[string]bool)
		for _, tool := range tools.Tools {
			nameSet[tool.Name] = true
		}
		for _, name := range expectedNames {
			t.Run("ListTools has "+name+" tool", func(t *testing.T) {
				if nameSet[name] != true {
					t.Fatalf("tool %s not found", name)
					return
				}
			})
		}
	})
}
