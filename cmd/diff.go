package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes",
	Long:  `Show difference between system dotfiles and target dotfiles.`,
	PreRun: func(cmd *cobra.Command, args []string ) {
		viper.BindPFlag("summaryView", cmd.PersistentFlags().Lookup("summary"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeDiff()
	},
}

func init() {
	diffCmd.PersistentFlags().BoolP("summary", "s", false, "Summary of the diff")
	rootCmd.AddCommand(diffCmd)
}

func executeDiff() error {
	config, err := internal.InitConfig()
	if err != nil {
		return err
	}

	internal.ShowDiffs(config)
	return nil
}
