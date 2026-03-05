package cmd

import (
	"fmt"

	"github.com/typemd/typemd/core"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new typemd vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := core.NewVault(".")
		if err := vault.Init(); err != nil {
			return err
		}
		fmt.Println("Initialized vault at .typemd/")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
