package cmd

import (
	"fmt"

	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Checks config file syntax and paths",
	Long:  `Verifies that the configuration file contains valid syntax and defined paths exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeVerify()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func executeVerify() error {
	_, err := internal.InitConfig()
	fmt.Printf("Using config file %s\n", viper.ConfigFileUsed())

	if err != nil {
		fmt.Println("Config file syntax is invalid.")
		return err
	} else {
		fmt.Println("Config file syntax is valid.")
		return nil
	}
}
