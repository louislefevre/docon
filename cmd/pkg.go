package cmd

import (
	"github.com/louislefevre/docon/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Generate a list of installed packages",
	Long:  `Retrieves a list of all packages installed on the system and sends it to a file in the target repository`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("dry", cmd.PersistentFlags().Lookup("dry-run"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return executePkg()
	},
}

func init() {
	pkgCmd.PersistentFlags().BoolP("dry-run", "d", false, "Run the command without actually changing anything")
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
