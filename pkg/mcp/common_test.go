package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/env"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/remote"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/store"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/versions"
	"sigs.k8s.io/controller-runtime/tools/setup-envtest/workflows"
	"testing"
)

// envTest has an expensive setup, so we only want to do it once per entire test run.
var envTest *envtest.Environment
var envTestRestConfig *rest.Config

func TestMain(m *testing.M) {
	// Set up
	envTestDir, err := store.DefaultStoreDir()
	if err != nil {
		panic(err)
	}
	envTestEnv := &env.Env{
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
	envTestEnv.CheckCoherence()
	workflows.Use{}.Do(envTestEnv)
	versionDir := envTestEnv.Platform.Platform.BaseName(*envTestEnv.Version.AsConcrete())
	envTest = &envtest.Environment{
		BinaryAssetsDirectory: filepath.Join(envTestDir, "k8s", versionDir),
	}
	envTestRestConfig, _ = envTest.Start()

	// Test!
	code := m.Run()

	// Tear down
	if envTest != nil {
		_ = envTest.Stop()
	}
	os.Exit(code)
}

type mcpContext struct {
	ctx        context.Context
	tempDir    string
	testServer *httptest.Server
	cancel     context.CancelFunc
	mcpClient  *client.SSEMCPClient
}

func (c *mcpContext) beforeEach(t *testing.T) {
	var err error
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.tempDir = t.TempDir()
	c.withKubeConfig(nil)
	c.testServer = server.NewTestServer(NewSever().server)
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

func testCase(t *testing.T, test func(c *mcpContext)) {
	mcpCtx := &mcpContext{}
	mcpCtx.beforeEach(t)
	defer mcpCtx.afterEach()
	test(mcpCtx)
}

func (c *mcpContext) withKubeConfig(rc *rest.Config) *api.Config {
	fakeConfig := api.NewConfig()
	fakeConfig.CurrentContext = "fake-context"
	fakeConfig.Contexts["fake-context"] = api.NewContext()
	fakeConfig.Contexts["fake-context"].Cluster = "fake"
	fakeConfig.Contexts["fake-context"].AuthInfo = "fake"
	fakeConfig.Clusters["fake"] = api.NewCluster()
	fakeConfig.Clusters["fake"].Server = "https://example.com"
	fakeConfig.AuthInfos["fake"] = api.NewAuthInfo()
	if rc != nil {
		fakeConfig.Clusters["fake"].Server = rc.Host
		fakeConfig.Clusters["fake"].CertificateAuthorityData = rc.TLSClientConfig.CAData
		fakeConfig.AuthInfos["fake"].ClientKeyData = rc.TLSClientConfig.KeyData
		fakeConfig.AuthInfos["fake"].ClientCertificateData = rc.TLSClientConfig.CertData
	}
	kubeConfig := filepath.Join(c.tempDir, "config")
	_ = clientcmd.WriteToFile(*fakeConfig, kubeConfig)
	_ = os.Setenv("KUBECONFIG", kubeConfig)
	return fakeConfig
}

func (c *mcpContext) withEnvTest() {
	c.withKubeConfig(envTestRestConfig)
}

func (c *mcpContext) newKubernetesClient() *kubernetes.Clientset {
	c.withEnvTest()
	pathOptions := clientcmd.NewDefaultPathOptions()
	cfg, _ := clientcmd.BuildConfigFromFlags("", pathOptions.GetDefaultFilename())
	kubernetesClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	return kubernetesClient
}

func (c *mcpContext) callTool(name string, args map[string]interface{}) (*mcp.CallToolResult, error) {
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Params.Name = name
	callToolRequest.Params.Arguments = args
	return c.mcpClient.CallTool(c.ctx, callToolRequest)
}
