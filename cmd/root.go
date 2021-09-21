package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "docon",
	Short: "A brief description of your application",
	Long:  `Docon is a command line tool used for maintaining Linux dotfiles.`,
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "",
		"config file (default is $HOME/.config/docon/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func ExecuteRoot() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}
