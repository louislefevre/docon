package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local dotfiles with repository",
	Long:  `Retrieves all dotfiles from the system and updates them in the repository.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("only", cmd.PersistentFlags().Lookup("only"))
		viper.BindPFlag("ignore", cmd.PersistentFlags().Lookup("ignore"))
		viper.BindPFlag("commit", cmd.PersistentFlags().Lookup("commit"))
		viper.BindPFlag("message", cmd.PersistentFlags().Lookup("message"))
		viper.BindPFlag("dryRun", cmd.PersistentFlags().Lookup("dry-run"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSync()
	},
}

func init() {
	syncCmd.PersistentFlags().StringSliceP("only", "o", []string{}, "Only sync these files (config override)")
	syncCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "Ignore these files (config override)")
	syncCmd.PersistentFlags().BoolP("commit", "c", false, "Commit dotfiles after syncing")
	syncCmd.PersistentFlags().StringP("message", "m", "", "Global commit message (config override)")
	syncCmd.PersistentFlags().BoolP("dry-run", "d", false, "Run the command without actually changing anything")
	rootCmd.AddCommand(syncCmd)
}

func executeSync() error {
	config, err := internal.InitConfig()
	if err != nil {
		return err
	}

	err = internal.SyncFiles(config)
	if err != nil {
		return err
	}

	if viper.GetBool("commit") {
		if err = internal.CommitAll(config); err != nil {
			return err
		}
	}

	return nil
}
