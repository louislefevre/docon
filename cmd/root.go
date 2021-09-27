package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "docon",
	Short: "A brief description of your application",
	Long:  `Docon is a command line tool used for maintaining Linux dotfiles.`,
	SilenceErrors: true,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/docon/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringSliceP("only", "o", []string{}, "Only use these groups (config override)")
	viper.BindPFlag("only", rootCmd.PersistentFlags().Lookup("only"))

	rootCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "Ignore these groups (config override)")
	viper.BindPFlag("ignore", rootCmd.PersistentFlags().Lookup("ignore"))
}

func ExecuteRoot() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
