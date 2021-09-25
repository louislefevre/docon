package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var filePaths []string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes",
	Long:  `Show difference between system dotfiles and target dotfiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePaths = args
		return executeDiff()
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func executeDiff() error {
	config, err := internal.InitConfig()
	if err != nil {
		return err
	}

	internal.ShowDiffs(config, filePaths)
	return nil
}
