package mcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/manusa/kubernetes-mcp-server/pkg/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Sever) initResources() {
	s.server.AddTool(mcp.NewTool(
		"resources_create_or_update",
		mcp.WithDescription("Create or update a Kubernetes resource in the current cluster by providing a YAML or JSON representation of the resource"),
		mcp.WithString("resource",
			mcp.Description("A JSON or YAML containing a representation of the Kubernetes resource. Should include top-level fields such as apiVersion,kind,metadata, and spec"),
			mcp.Required(),
		),
	), resourcesCreateOrUpdate)
}

func resourcesCreateOrUpdate(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to create or update resources: %v", err)), nil
	}
	resource := ctr.Params.Arguments["resource"]
	if resource == nil || resource == "" {
		return NewTextResult("", errors.New("failed to create or update resources, missing argument resource")), nil
	}
	ret, err := k.ResourcesCreateOrUpdate(ctx, resource.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to create or update resources: %v", err)), nil
	}
	return NewTextResult(ret, err), nil
}
