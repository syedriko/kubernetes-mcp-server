package mcp

import (
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Sever struct {
	server *server.MCPServer
}

func NewSever() *Sever {
	s := &Sever{
		server: server.NewMCPServer(
			version.BinaryName,
			version.Version,
			server.WithResourceCapabilities(true, true),
			server.WithPromptCapabilities(true),
			server.WithLogging(),
		),
	}
	s.initConfiguration()
	s.initPods()
	s.initResources()
	return s
}

func (s *Sever) ServeStdio() error {
	return server.ServeStdio(s.server)
}

func NewTextResult(content string, err error) *mcp.CallToolResult {
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []interface{}{
				mcp.TextContent{
					Type: "text",
					Text: err.Error(),
				},
			},
		}
	}
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: content,
			},
		},
	}
}
