package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local dotfiles with repository",
	Long:  `Retrieves all dotfiles from the system and updates them in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ExecuteSync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&internal.RepoPath, "repo", "r", "", "path to repository directory")
	syncCmd.MarkFlagRequired("repo")
}
