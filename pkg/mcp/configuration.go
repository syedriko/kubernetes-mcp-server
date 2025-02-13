package mcp

import (
	"context"
	"github.com/manusa/kubernetes-mcp-server/pkg/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
)

func configurationView(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cfg, err := kubernetes.ConfigurationView()
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: cfg,
			},
		},
	}, nil
}
