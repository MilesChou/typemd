package cmd

import (
	"github.com/MilesChou/typemd/core"
	mcpserver "github.com/MilesChou/typemd/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server via stdio",
	Long:  `Start a Model Context Protocol (MCP) server using stdio transport. This allows AI clients to query the vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := core.NewVault(vaultPath)
		if vaultPath == "" {
			vault = core.NewVault(".")
		}
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		return mcpserver.Serve(vault)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
