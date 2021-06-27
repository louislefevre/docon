package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Checks config file syntax and paths",
	Long:  `Verifies that the configuration file contains valid syntax and defined paths exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ExecuteVerify()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
