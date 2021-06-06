package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Process dotfiles",
	Long:  `Retrieves all dotfiles from the system and updates them in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeRun()
	},
}

var configFile string

type Configuration map[string]ConfigGroup

type ConfigGroup struct {
	Path    string   `mapstructure:"path"`
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config/docon/config.yaml)")
}

func executeRun() {
	cobra.CheckErr(initConfig())
}

func initConfig() error {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s/.config/docon/", home))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("missing config file: %s", err)
		} else {
			return fmt.Errorf("an error occurred: %s", err)
		}
	}

	config := make(Configuration)
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("cannot parse config file: %s", err)
	}

	for name, group := range config {
		if group.Path == "" {
			return fmt.Errorf("cannot parse config file: %s has no defined path", name)
		} else if _, err := os.Stat(group.Path); os.IsNotExist(err) {
			return fmt.Errorf("cannot parse config file: %s does not exist", group.Path)
		}
		fmt.Println(group.Path)
	}

	return nil
}
