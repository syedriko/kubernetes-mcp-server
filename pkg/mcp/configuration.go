package mcp

import (
	"context"
	"fmt"
	"github.com/manusa/kubernetes-mcp-server/pkg/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initConfiguration() []server.ServerTool {
	return []server.ServerTool{
		{mcp.NewTool("configuration_view",
			mcp.WithDescription("Get the current Kubernetes configuration content as a kubeconfig YAML"),
		), configurationView},
	}
}

func configurationView(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ret, err := kubernetes.ConfigurationView()
	if err != nil {
		err = fmt.Errorf("failed to get configuration: %v", err)
	}
	return NewTextResult(ret, err), nil
}
