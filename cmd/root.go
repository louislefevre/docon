package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docon",
	Short: "A brief description of your application",
	Long:  `Docon is a command line tool used for maintaining Linux dotfiles.`,
}

func init() {}

func ExecuteRoot() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}
