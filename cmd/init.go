package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new typemd vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Init(); err != nil {
			return err
		}
		fmt.Printf("Initialized vault at %s\n", vault.Dir())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
