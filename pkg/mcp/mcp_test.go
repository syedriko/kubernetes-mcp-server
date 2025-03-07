package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"strings"
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

func TestToolsInOpenShift(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		defer c.inOpenShift()() // n.b. two sets of parentheses to invoke the first function
		c.mcpServer.server.AddTools(c.mcpServer.initResources()...)
		tools, err := c.mcpClient.ListTools(c.ctx, mcp.ListToolsRequest{})
		t.Run("ListTools returns tools", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call ListTools failed %v", err)
				return
			}
		})
		t.Run("ListTools has resources_list tool with OpenShift hint", func(t *testing.T) {
			if tools.Tools[10].Name != "resources_list" {
				t.Fatalf("tool resources_list not found")
				return
			}
			if !strings.Contains(tools.Tools[10].Description, ", route.openshift.io/v1 Route") {
				t.Fatalf("tool resources_list does not have OpenShift hint, got %s", tools.Tools[9].Description)
				return
			}
		})
	})

}
