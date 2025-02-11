package mcp

import (
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	"github.com/mark3labs/mcp-go/server"
)

func Start() {
	s := server.NewMCPServer(
		version.BinaryName,
		version.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}
