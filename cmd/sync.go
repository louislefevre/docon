package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var commitFiles bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local dotfiles with repository",
	Long:  `Retrieves all dotfiles from the system and updates them in the repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeSync()
	},
}

func init() {
	syncCmd.PersistentFlags().BoolVarP(&commitFiles, "commit", "c", false, "Commit dotfiles after syncing")
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

	if commitFiles {
		if err = internal.CommitAll(config); err != nil {
			return err
		}
	}

	return nil
}
