package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var RepoPath string

type configuration struct {
	repoPath string
	mapping  configMap
}

type configMap map[string]configGroup

type configGroup struct {
	path     string   `mapstructure:"path"`
	included []string `mapstructure:"include"`
	excluded []string `mapstructure:"exclude"`
}

func initConfig() (configuration, error) {
	var config configuration

	if home, err := homedir.Dir(); err == nil {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/docon/config.yaml", home))
	} else {
		return config, fmt.Errorf("failed to find home directory\n%s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to load config file\n%s", err)
	}

	mapping := make(configMap)
	if err := viper.Unmarshal(&mapping); err != nil {
		return config, fmt.Errorf("failed to parse config file\n%s", err)
	}

	for name, group := range mapping {
		if group.path == "" {
			return config, fmt.Errorf("failed to parse config file\n%s: no defined path", name)
		} else if fileInfo, err := os.Stat(group.path); os.IsNotExist(err) {
			return config, fmt.Errorf("failed to parse config file\n%s", err)
		} else if !fileInfo.IsDir() {
			return config, fmt.Errorf("failed to parse config file\n%s: is not a directory", group.path)
		}

		for i, file := range group.included {
			filePath := filepath.Join(group.path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return config, fmt.Errorf("failed to parse config file\n%s", err)
			}
			group.included[i] = filePath
		}

		for i, file := range group.excluded {
			filePath := filepath.Join(group.path, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return config, fmt.Errorf("failed to parse config file\n%s", err)
			}
			group.excluded[i] = filePath
		}
	}

	if _, err := os.Stat(RepoPath); os.IsNotExist(err) {
		return config, fmt.Errorf("failed to find repo directory\n%s", err)
	}

	config.repoPath = RepoPath
	config.mapping = mapping

	return config, nil
}
