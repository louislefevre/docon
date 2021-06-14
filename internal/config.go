package internal

import (
	"fmt"
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
	Path     string   `mapstructure:"path"`
	Included []string `mapstructure:"include"`
	Excluded []string `mapstructure:"exclude"`
}

func initConfig() (configuration, error) {
	var config configuration

	mapping, err := loadConfig()
	if err != nil {
		return config, err
	}

	config.repoPath = RepoPath
	config.mapping = mapping

	return config, nil
}

func loadConfig() (configMap, error) {
	if home, err := homedir.Dir(); err == nil {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/docon/config.yaml", home))
	} else {
		return nil, fmt.Errorf("failed to find home directory\n%s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config file\n%s", err)
	}

	mapping := make(configMap)
	if err := viper.Unmarshal(&mapping); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file\n%s", err)
	}

	if err := parseConfig(mapping); err != nil {
		return nil, fmt.Errorf("failed to parse config file\n%s", err)
	}

	return mapping, nil
}

func parseConfig(mapping configMap) error {
	for name, group := range mapping {
		if group.Path == "" {
			return fmt.Errorf("%s has no defined path", name)
		}

		for i, file := range group.Included {
			group.Included[i] = filepath.Join(group.Path, file)
		}

		for i, file := range group.Excluded {
			group.Excluded[i] = filepath.Join(group.Path, file)
		}
	}

	return nil
}
