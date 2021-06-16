package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes.",
	Long:  `Show difference between system dotfiles and target dotfiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ExecuteDiff(args)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVarP(&internal.RepoPath, "repo", "r", "", "path to repository directory")
	diffCmd.MarkFlagRequired("repo")
}
