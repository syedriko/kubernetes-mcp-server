package mcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initPods() []server.ServerTool {
	return []server.ServerTool{
		{mcp.NewTool("pods_list",
			mcp.WithDescription("List all the Kubernetes pods in the current cluster from all namespaces"),
		), s.podsListInAllNamespaces},
		{mcp.NewTool("pods_list_in_namespace",
			mcp.WithDescription("List all the Kubernetes pods in the specified namespace in the current cluster"),
			mcp.WithString("namespace", mcp.Description("Namespace to list pods from"), mcp.Required()),
		), s.podsListInNamespace},
		{mcp.NewTool("pods_get",
			mcp.WithDescription("Get a Kubernetes Pod in the current or provided namespace with the provided name"),
			mcp.WithString("namespace", mcp.Description("Namespace to get the Pod from")),
			mcp.WithString("name", mcp.Description("Name of the Pod"), mcp.Required()),
		), s.podsGet},
		{mcp.NewTool("pods_delete",
			mcp.WithDescription("Delete a Kubernetes Pod in the current or provided namespace with the provided name"),
			mcp.WithString("namespace", mcp.Description("Namespace to delete the Pod from")),
			mcp.WithString("name", mcp.Description("Name of the Pod to delete"), mcp.Required()),
		), s.podsDelete},
		{mcp.NewTool("pods_log",
			mcp.WithDescription("Get the logs of a Kubernetes Pod in the current or provided namespace with the provided name"),
			mcp.WithString("namespace", mcp.Description("Namespace to get the Pod logs from")),
			mcp.WithString("name", mcp.Description("Name of the Pod to get the logs from"), mcp.Required()),
		), s.podsLog},
		{mcp.NewTool("pods_run",
			mcp.WithDescription("Run a Kubernetes Pod in the current or provided namespace with the provided container image and optional name"),
			mcp.WithString("namespace", mcp.Description("Namespace to run the Pod in")),
			mcp.WithString("name", mcp.Description("Name of the Pod (Optional, random name if not provided)")),
			mcp.WithString("image", mcp.Description("Container Image to run in the Pod"), mcp.Required()),
			mcp.WithNumber("port", mcp.Description("TCP/IP port to expose from the Pod container (Optional, no port exposed if not provided)")),
		), s.podsRun},
	}
}

func (s *Server) podsListInAllNamespaces(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ret, err := s.k.PodsListInAllNamespaces(ctx)
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in all namespaces: %v", err)), nil
	}
	return NewTextResult(ret, err), nil
}

func (s *Server) podsListInNamespace(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		return NewTextResult("", errors.New("failed to list pods in namespace, missing argument namespace")), nil
	}
	ret, err := s.k.PodsListInNamespace(ctx, ns.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to list pods in namespace %s: %v", ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}

func (s *Server) podsGet(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		return NewTextResult("", errors.New("failed to get pod, missing argument name")), nil
	}
	ret, err := s.k.PodsGet(ctx, ns.(string), name.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}

func (s *Server) podsDelete(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		return NewTextResult("", errors.New("failed to delete pod, missing argument name")), nil
	}
	ret, err := s.k.PodsDelete(ctx, ns.(string), name.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to delete pod %s in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}

func (s *Server) podsLog(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.Params.Arguments["namespace"]
	if ns == nil {
		ns = ""
	}
	name := ctr.Params.Arguments["name"]
	if name == nil {
		return NewTextResult("", errors.New("failed to get pod log, missing argument name")), nil
	}
	ret, err := s.k.PodsLog(ctx, ns.(string), name.(string))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s log in namespace %s: %v", name, ns, err)), nil
	} else if ret == "" {
		ret = fmt.Sprintf("The pod %s in namespace %s has not logged any message yet", name, ns)
	}
	return NewTextResult(ret, err), nil
}

func (s *Server) podsRun(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	ret, err := s.k.PodsRun(ctx, ns.(string), name.(string), image.(string), int32(port.(float64)))
	if err != nil {
		return NewTextResult("", fmt.Errorf("failed to get pod %s log in namespace %s: %v", name, ns, err)), nil
	}
	return NewTextResult(ret, err), nil
}
