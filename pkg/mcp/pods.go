package mcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/manusa/kubernetes-mcp-server/pkg/kubernetes"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Sever) initPods() {
	s.server.AddTool(mcp.NewTool(
		"pods_list",
		mcp.WithDescription("List all the Kubernetes pods in the current cluster from all namespaces"),
	), podsListInAllNamespaces)
	s.server.AddTool(mcp.NewTool(
		"pods_list_in_namespace",
		mcp.WithDescription("List all the Kubernetes pods in the specified namespace in the current cluster"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to list pods from"),
			mcp.Required(),
		),
	), podsListInNamespace)
}

func podsListInAllNamespaces(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in all namespaces: %v", err)), nil
	}
	ret, err := k.PodsListInAllNamespaces(ctx)
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in all namespaces: %v", err)), nil
	}
	return NewTextResult(ret, err), nil
}

func podsListInNamespace(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in namespace: %v", err)), nil
	}
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		return NewTextResult("", errors.New("failed to list pods in namespace, missing argument namespace")), nil
	}
	ret, err := k.PodsListInNamespace(ctx, ns.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in namespace %s: %v", ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}
