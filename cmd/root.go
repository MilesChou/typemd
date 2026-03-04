package cmd

import (
	"os"

	"github.com/MilesChou/typemd/tui"
	"github.com/spf13/cobra"
)

var vaultPath string

var rootCmd = &cobra.Command{
	Use:   "tmd",
	Short: "A local-first CLI knowledge management tool",
	// 不帶子指令時啟動 TUI
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start(vaultPath)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultPath, "vault", "", "path to vault directory (default: current directory)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
