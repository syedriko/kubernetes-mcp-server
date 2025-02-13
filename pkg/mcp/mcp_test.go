package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"testing"
)

func TestTools(t *testing.T) {
	t.Run("Has configuration_view tool", testCase(func(t *testing.T, c *mcpContext) {
		tools, err := c.mcpClient.ListTools(c.ctx, mcp.ListToolsRequest{})
		if tools.Tools[0].Name != "configuration_view" {
			t.Fatalf("tool name mismatch %v", err)
			return
		}
	}))
}
