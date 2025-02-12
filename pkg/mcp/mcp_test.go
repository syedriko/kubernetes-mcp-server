package mcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"testing"
)

func TestCapabilities(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewSever()
	testServer := server.NewTestServer(s.server)
	defer testServer.Close()
	c, err := client.NewSSEMCPClient(testServer.URL + "/sse")
	defer func() {
		_ = c.Close()
	}()
	if err != nil {
		t.Fatal(err)
		return
	}
	if err = c.Start(ctx); err != nil {
		t.Fatal(err)
		return
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.33.7"}
	ir, err := c.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Print(ir)
	cancel()
}
