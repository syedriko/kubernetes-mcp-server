#!/usr/bin/env node

const path = require('path')
const childProcess = require('child_process');

const BINARY_MAP = {
  darwin_x64: {name: 'kubernetes-mcp-server-darwin-amd64', suffix: ''},
  darwin_arm64: {name: 'kubernetes-mcp-server-darwin-arm64', suffix: ''},
  linux_x64: {name: 'kubernetes-mcp-server-linux-amd64', suffix: ''},
  linux_arm64: {name: 'kubernetes-mcp-server-linux-arm64', suffix: ''},
  win32_x64: {name: 'kubernetes-mcp-server-windows-amd64', suffix: '.exe'},
  win32_arm64: {name: 'kubernetes-mcp-server-windows-arm64', suffix: '.exe'},
};

const binary = BINARY_MAP[`${process.platform}_${process.arch}`];

const resolveBinaryPath = () => {
  try {
    // Resolving will fail if the optionalDependency was not installed
    return require.resolve(`${binary.name}/bin/${binary.name}${binary.suffix}`)
  } catch (e) {
    return path.join(__dirname, '..', `${binary.name}${binary.suffix}`)
  }
};

childProcess.execFileSync(resolveBinaryPath(), process.argv.slice(2), {
  stdio: 'inherit',
});

