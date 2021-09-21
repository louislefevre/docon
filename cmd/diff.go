package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes.",
	Long:  `Show difference between system dotfiles and target dotfiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.ExecuteDiff(args)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
