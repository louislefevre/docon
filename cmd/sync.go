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
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSync()
	},
}

func init() {
	syncCmd.PersistentFlags().BoolP("commit", "c", false, "Commit dotfiles after syncing")
	viper.BindPFlag("commit", syncCmd.PersistentFlags().Lookup("commit"))

	syncCmd.PersistentFlags().StringP("message", "m", "", "Global commit message (config override)")
	viper.BindPFlag("message", syncCmd.PersistentFlags().Lookup("message"))

	syncCmd.PersistentFlags().BoolP("dry-run", "d", false, "Run the command without actually changing anything")
	viper.BindPFlag("dry", syncCmd.PersistentFlags().Lookup("dry-run"))

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
