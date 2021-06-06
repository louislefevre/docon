package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

type Configuration map[string]ConfigGroup

type ConfigGroup struct {
	Path    string   `mapstructure:"path"`
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func executeRun() {
	config, err := initConfig()
	cobra.CheckErr(err)
	fmt.Println(config)
}

func initConfig() (Configuration, error) {
	if home, err := homedir.Dir(); err == nil {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/docon/config.yaml", home))
	} else {
		return nil, fmt.Errorf("failed to find home directory\n%s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config file\n%s", err)
	}

	config := make(Configuration)
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file\n%s", err)
	}

	for name, group := range config {
		if group.Path == "" {
			return nil, fmt.Errorf("failed to parse config file\n%s: no defined path", name)
		} else if _, err := os.Stat(group.Path); os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to parse config file\n%s", err)
		}

		for _, file := range group.Include {
			filePath := filepath.Join(group.Path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to parse config file\n%s", err)
			}
		}

		for _, file := range group.Exclude {
			filePath := filepath.Join(group.Path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to parse config file\n%s", err)
			}
		}
	}

	return config, nil
}
