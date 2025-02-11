package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubernetes-mcp-server [command] [options]",
	Short: "Kubernetes Model Context Protocol (MCP) server",
	Long: `
Kubernetes Model Context Protocol (MCP) server

  # show this help
  kubernetes-mcp-server -h

  # TODO: add more examples`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
