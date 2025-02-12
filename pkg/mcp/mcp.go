package mcp

import (
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	"github.com/mark3labs/mcp-go/server"
)

type Sever struct {
	server *server.MCPServer
}

func NewSever() *Sever {
	s := server.NewMCPServer(
		version.BinaryName,
		version.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	return &Sever{s}
}

func (s *Sever) ServeStdio() error {
	return server.ServeStdio(s.server)
}
