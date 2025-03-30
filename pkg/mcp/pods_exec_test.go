package mcp

import (
	"bytes"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"testing"
)

func TestPodsExec(t *testing.T) {
	testCase(t, func(c *mcpContext) {
		mockServer := NewMockServer()
		defer mockServer.Close()
		c.withKubeConfig(mockServer.config)
		mockServer.Handle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/v1/namespaces/default/pods/pod-to-exec/exec" {
				return
			}
			var stdin, stdout bytes.Buffer
			ctx, err := createHTTPStreams(w, req, &StreamOptions{
				Stdin:  &stdin,
				Stdout: &stdout,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			defer ctx.conn.Close()
			_, _ = io.WriteString(ctx.stdoutStream, strings.Join(req.URL.Query()["command"], " "))
			_, _ = io.WriteString(ctx.stdoutStream, "\ntotal 0\n")
		}))
		mockServer.Handle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/api/v1/namespaces/default/pods/pod-to-exec" {
				return
			}
			writeObject(w, &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "pod-to-exec",
				},
				Spec: v1.PodSpec{Containers: []v1.Container{{Name: "container-to-exec"}}},
			})
		}))
		toolResult, err := c.callTool("pods_exec", map[string]interface{}{
			"namespace": "default",
			"name":      "pod-to-exec",
			"command":   []interface{}{"ls", "-l"},
		})
		t.Run("pods_exec returns command output", func(t *testing.T) {
			if err != nil {
				t.Fatalf("call tool failed %v", err)
			}
			if toolResult.IsError {
				t.Fatalf("call tool failed")
			}
			if toolResult.Content[0].(mcp.TextContent).Text != "ls -l\ntotal 0\n" {
				t.Errorf("unexpected result %v", toolResult.Content[0].(mcp.TextContent).Text)
			}
		})

	})
}
