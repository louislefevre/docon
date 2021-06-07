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

var RepoPath string

type Configuration struct {
	RepoPath string
	ConfigMap ConfigMap
}

type ConfigMap map[string]ConfigGroup

type ConfigGroup struct {
	Path    string   `mapstructure:"path"`
	Included []string `mapstructure:"include"`
	Excluded []string `mapstructure:"exclude"`
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&RepoPath, "repo", "r", "", "path to repository directory")
	runCmd.MarkFlagRequired("repo")
}

func executeRun() {
	config, err := initConfig()
	cobra.CheckErr(err)
	err = processConfiguration(config)
	cobra.CheckErr(err)
}

func initConfig() (Configuration, error) {
	var config Configuration

	if home, err := homedir.Dir(); err == nil {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/docon/config.yaml", home))
	} else {
		return config, fmt.Errorf("failed to find home directory\n%s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to load config file\n%s", err)
	}

	configMap := make(ConfigMap)
	if err := viper.Unmarshal(&configMap); err != nil {
		return config, fmt.Errorf("failed to parse config file\n%s", err)
	}

	for name, group := range configMap {
		if group.Path == "" {
			return config, fmt.Errorf("failed to parse config file\n%s: no defined path", name)
		} else if _, err := os.Stat(group.Path); os.IsNotExist(err) {
			return config, fmt.Errorf("failed to parse config file\n%s", err)
		}

		for i, file := range group.Included {
			filePath := filepath.Join(group.Path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return config, fmt.Errorf("failed to parse config file\n%s", err)
			}
			group.Included[i] = filePath
		}

		for i, file := range group.Excluded {
			filePath := filepath.Join(group.Path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return config, fmt.Errorf("failed to parse config file\n%s", err)
			}
			group.Excluded[i] = filePath
		}
	}

	if _, err := os.Stat(RepoPath); os.IsNotExist(err) {
		return config, fmt.Errorf("failed to find repo directory\n%s", err)
	}

	config.RepoPath = RepoPath
	config.ConfigMap = configMap

	return config, nil
}
