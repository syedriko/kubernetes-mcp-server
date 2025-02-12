package cmd

import (
	"fmt"
	"github.com/manusa/kubernetes-mcp-server/pkg/mcp"
	"github.com/manusa/kubernetes-mcp-server/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kubernetes-mcp-server [command] [options]",
	Short: "Kubernetes Model Context Protocol (MCP) server",
	Long: `
Kubernetes Model Context Protocol (MCP) server

  # show this help
  kubernetes-mcp-server -h

  # shows version information
  kubernetes-mcp-server --version

  # TODO: add more examples`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version") {
			fmt.Println(version.Version)
			return
		}
		if err := mcp.NewSever().ServeStdio(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version information and quit")
	_ = viper.BindPFlags(rootCmd.Flags())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
