# Kubernetes MCP Server

<p align="center">
  <a href="https://github.com/manusa/kubernetes-mcp-server/blob/main/LICENSE">
    <img alt="GitHub License" src="https://img.shields.io/github/license/manusa/kubernetes-mcp-server" /></a>
  <a href="https://www.npmjs.com/package/kubernetes-mcp-server">
    <img alt="npm" src="https://img.shields.io/npm/v/kubernetes-mcp-server" /></a>
  <a href="https://github.com/manusa/kubernetes-mcp-server/releases/latest">
    <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/manusa/kubernetes-mcp-server?sort=semver" /></a>
  <a href="https://github.com/manusa/kubernetes-mcp-server/actions/workflows/build.yaml">
    <img src="https://github.com/manusa/kubernetes-mcp-server/actions/workflows/build.yaml/badge.svg" alt="Build status badge" /></a>
</p>


<p align="center">
  <a href="#features">✨ Features</a> |
  <a href="#getting-started">🚀 Getting Started</a>
</p>

## ✨ Features <a id="features"></a>

A powerful and flexible Kubernetes MCP server implementation with support for OpenShift.

- **✅ Configuration**: View and manage the Kubernetes `.kube/config`.
  - **View** the current configuration.
- **✅ Generic Kubernetes Resources**: Perform operations on any Kubernetes resource.
  - Any CRUD operation (Create or Update, Get, List, Delete).
- **✅ Pods**: Perform Pod-specific operations.
  - **List** pods in all namespaces or in a specific namespace.
  - **Get** a pod by name from the specified namespace.
  - **Delete** a pod by name from the specified namespace.
  - **Show logs** for a pod by name from the specified namespace.
  - **Run** a container image in a pod and optionally expose it.

## 🚀 Getting Started <a id="getting-started"></a>

### Claude Desktop

#### Using npx

If you have npm installed, this is the fastest way to get started with `kubernetes-mcp-server`.

Open your `claude_desktop_config.json` and add the mcp server to the list of `mcpServers`:
``` json
{
  "mcpServers": {
    "kubernetes-mcp": {
      "command": "npx",
      "args": ["-y", "kubernetes-mcp-server@latest"]
  }
}
```

