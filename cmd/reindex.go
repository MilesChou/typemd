package cmd

import (
	"fmt"

	"github.com/typemd/typemd/core"
	"github.com/spf13/cobra"
)

var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Sync objects from disk and rebuild search index",
	Long:  `Scan the objects directory, sync all files to the database, and rebuild the FTS5 search index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := core.NewVault(vaultPath)
		if vaultPath == "" {
			vault = core.NewVault(".")
		}
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		if err := vault.SyncIndex(); err != nil {
			return err
		}

		fmt.Println("Index synced successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)
}
