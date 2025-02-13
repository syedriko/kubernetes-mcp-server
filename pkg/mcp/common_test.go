package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/afero"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/env"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/remote"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/store"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/versions"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/workflows"
)

func setupEnvTest() *envtest.Environment {
	envTestDir, err := store.DefaultStoreDir()
	if err != nil {
		panic(err)
	}
	envTest := &env.Env{
		FS:  afero.Afero{Fs: afero.NewOsFs()},
		Out: os.Stdout,
		Client: &remote.HTTPClient{
			IndexURL: remote.DefaultIndexURL,
		},
		Platform: versions.PlatformItem{
			Platform: versions.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
		},
		Version: versions.AnyVersion,
		Store:   store.NewAt(envTestDir),
	}
	envTest.CheckCoherence()
	workflows.Use{}.Do(envTest)
	versionDir := envTest.Platform.Platform.BaseName(*envTest.Version.AsConcrete())
	return &envtest.Environment{
		BinaryAssetsDirectory: filepath.Join(envTestDir, "k8s", versionDir),
	}
}

func withKubeConfig(t *testing.T, c *rest.Config) *api.Config {
	fakeConfig := api.NewConfig()
	fakeConfig.CurrentContext = "fake-context"
	fakeConfig.Clusters["fake"] = api.NewCluster()
	fakeConfig.Clusters["fake"].Server = c.Host
	fakeConfig.Clusters["fake"].CertificateAuthorityData = c.TLSClientConfig.CAData
	fakeConfig.Contexts["fake-context"] = api.NewContext()
	fakeConfig.Contexts["fake-context"].Cluster = "fake"
	fakeConfig.Contexts["fake-context"].AuthInfo = "fake"
	fakeConfig.AuthInfos["fake"] = api.NewAuthInfo()
	fakeConfig.AuthInfos["fake"].ClientKeyData = c.TLSClientConfig.KeyData
	fakeConfig.AuthInfos["fake"].ClientCertificateData = c.TLSClientConfig.CertData
	dir := t.TempDir()
	kubeConfig := filepath.Join(dir, "config")
	clientcmd.WriteToFile(*fakeConfig, kubeConfig)
	os.Setenv("KUBECONFIG", kubeConfig)
	return fakeConfig
}

type mcpContext struct {
	ctx        context.Context
	testServer *httptest.Server
	cancel     context.CancelFunc
	mcpClient  *client.SSEMCPClient
}

func (c *mcpContext) beforeEach(t *testing.T) {
	var err error
	c.testServer = server.NewTestServer(NewSever().server)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	if c.mcpClient, err = client.NewSSEMCPClient(c.testServer.URL + "/sse"); err != nil {
		t.Fatal(err)
		return
	}
	if err = c.mcpClient.Start(c.ctx); err != nil {
		t.Fatal(err)
		return
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.33.7"}
	_, err = c.mcpClient.Initialize(c.ctx, initRequest)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func (c *mcpContext) afterEach() {
	c.cancel()
	_ = c.mcpClient.Close()
	c.testServer.Close()
}

func testCase(test func(t *testing.T, c *mcpContext)) func(*testing.T) {
	return func(t *testing.T) {
		mcpCtx := &mcpContext{}
		mcpCtx.beforeEach(t)
		defer mcpCtx.afterEach()
		test(t, mcpCtx)
	}
}
