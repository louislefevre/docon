package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes",
	Long:  `Automatically commit target dotfiles to Git repository with pre-defined commit messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ExecuteCommit()
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
