package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docon",
	Short: "A brief description of your application",
	Long:  `Docon is a command line tool used for maintaining Linux dotfiles.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config/docon/config.yaml)")
}

func ExecuteRoot() {
	cobra.CheckErr(rootCmd.Execute())
}
