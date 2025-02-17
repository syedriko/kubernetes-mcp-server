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
	s.server.AddTool(mcp.NewTool(
		"pods_get",
		mcp.WithDescription("Get a Kubernetes Pod in the current or provided namespace with the provided name"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get the Pod from"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the Pod"),
			mcp.Required(),
		),
	), podsGet)
	s.server.AddTool(mcp.NewTool(
		"pods_log",
		mcp.WithDescription("Get the logs of a Kubernetes Pod in the current or provided namespace with the provided name"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to get the Pod logs from"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the Pod"),
			mcp.Required(),
		),
	), podsLog)
	s.server.AddTool(mcp.NewTool(
		"pods_run",
		mcp.WithDescription("Run a Kubernetes Pod in the current or provided namespace with the provided container image and optional name"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to run the Pod in"),
		),
		mcp.WithString("name",
			mcp.Description("Name of the Pod (Optional, random name if not provided)"),
		),
		mcp.WithString("image",
			mcp.Description("Container Image to run in the Pod"),
			mcp.Required(),
		),
		mcp.WithNumber("port",
			mcp.Description("TCP/IP port to expose from the Pod container (Optional, no port exposed if not provided)"),
		),
	), podsRun)
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

func podsGet(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod: %v", err)), nil
	}
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		return NewTextResult("", errors.New("failed to get pod, missing argument name")), nil
	}
	ret, err := k.PodsGet(ctx, ns.(string), name.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}

func podsLog(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod log: %v", err)), nil
	}
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		return NewTextResult("", errors.New("failed to get pod log, missing argument name")), nil
	}
	ret, err := k.PodsLog(ctx, ns.(string), name.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s log in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}

func podsRun(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	k, err := kubernetes.NewKubernetes()
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to run pod: %v", err)), nil
	}
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		name = ""
	}
	image := ctr.Params.Arguments["image"]
	if image == nil {
		return NewTextResult("", errors.New("failed to run pod, missing argument image")), nil
	}
	port := ctr.Params.Arguments["port"]
	if port == nil {
		port = float64(0)
	}
	ret, err := k.PodsRun(ctx, ns.(string), name.(string), image.(string), int32(port.(float64)))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s log in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}
