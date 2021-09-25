package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Generate a list of installed packages",
	Long:  `Retrieves a list of all packages installed on the system and sends it to a file in the target repository`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executePkg()
	},
}

func init() {
	rootCmd.AddCommand(pkgCmd)
}

func executePkg() error {
	config, err := internal.InitConfig()
	if err != nil {
		return err
	}

	err = internal.GenPackageList(config)
	if err != nil {
		return err
	}

	return nil
}
