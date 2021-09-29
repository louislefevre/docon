package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes",
	Long:  `Automatically commit target dotfiles to Git repository with pre-defined commit messages.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("only", cmd.PersistentFlags().Lookup("only"))
		viper.BindPFlag("ignore", cmd.PersistentFlags().Lookup("ignore"))
		viper.BindPFlag("message", cmd.PersistentFlags().Lookup("message"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeCommit()
	},
}

func init() {
	commitCmd.PersistentFlags().StringSliceP("only", "o", []string{}, "Only commit these files (config override)")
	commitCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "Ignore these files (config override)")
	commitCmd.PersistentFlags().StringP("message", "m", "", "Global commit message (config override)")
	rootCmd.AddCommand(commitCmd)
}

func executeCommit() error {
	config, err := internal.InitConfig()
	if err != nil {
		return err
	}

	err = internal.CommitAll(config)
	if err != nil {
		return err
	}

	return nil
}
