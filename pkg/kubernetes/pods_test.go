package kubernetes

import (
	"bytes"
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

func TestPodsExec(t *testing.T) {
	mockServer := NewMockServer()
	defer mockServer.Close()
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
		_, _ = io.WriteString(ctx.stdoutStream, "total 0\n")
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
	k8s := mockServer.NewKubernetes()
	out, err := k8s.PodsExec(context.Background(), "default", "pod-to-exec", "", []string{"ls", "-l"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "total 0\n" {
		t.Fatalf("unexpected output: %s", out)
	}
}
