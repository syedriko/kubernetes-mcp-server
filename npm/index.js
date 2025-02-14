#!/usr/bin/env node

const childProcess = require("child_process");

const BINARY_MAP = {
  darwin_x86: {name: "kubernetes-mcp-server-darwin-amd64", suffix: ''},
  darwin_arm64: {name: "kubernetes-mcp-server-darwin-arm64", suffix: ''},
  linux_x86: {name: "kubernetes-mcp-server-linux-amd64", suffix: ''},
  linux_arm64: {name: "kubernetes-mcp-server-linux-arm64", suffix: ''},
  win32_x86: {name: "kubernetes-mcp-server-windows-amd64", suffix: '.exe'},
  win32_arm64: {name: "kubernetes-mcp-server-windows-arm64", suffix: '.exe'},
};

const binary = BINARY_MAP[`${process.platform}_${process.arch}`];

module.exports.runBinary = function (...args) {
  // Resolving will fail if the optionalDependency was not installed
  childProcess.execFileSync(require.resolve(`${binary.name}/bin/${binary.name}+${binary.suffix}`), args, {
    stdio: "inherit",
  });
};
