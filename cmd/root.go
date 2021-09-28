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
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
		viper.BindPFlag("only", cmd.PersistentFlags().Lookup("only"))
		viper.BindPFlag("ignore", cmd.PersistentFlags().Lookup("ignore"))
	},
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/docon/config.yaml)")
	rootCmd.PersistentFlags().StringSliceP("only", "o", []string{}, "Only use these groups (config override)")
	rootCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "Ignore these groups (config override)")
}

func ExecuteRoot() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
