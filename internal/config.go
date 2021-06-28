package internal

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type configuration struct {
	TargetPath  string `mapstructure:"target"`
	PkglistPath string `mapstructure:"pkglist"`
	Git         struct {
		Dir       bool   `mapstructure:"dir"`
		CommitMsg string `mapstructure:"msg"`
		Author    struct {
			name  string `mapstructure:"name"`
			email string `mapstructure:"email"`
		}
	} `mapstructure:"git"`
	Sources map[string]struct {
		Path      string   `mapstructure:"path"`
		CommitMsg string   `mapstructure:"msg"`
		Included  []string `mapstructure:"include"`
		Excluded  []string `mapstructure:"exclude"`
		dotfiles  dotfiles
	} `mapstructure:"sources"`
	allDotfiles dotfiles
}

func initConfig() (*configuration, error) {
	var config configuration

	if err := loadConfig(&config); err != nil {
		return &config, err
	}

	if err := parseConfig(&config); err != nil {
		return &config, err
	}

	if err := verifyConfig(&config); err != nil {
		return &config, err
	}

	if err := gatherDotfiles(&config); err != nil {
		return &config, err
	}

	return &config, nil
}

func loadConfig(config *configuration) error {
	if home, err := homedir.Dir(); err == nil {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/docon/config.yaml", home))
	} else {
		return newError(err, "Failed to find home directory")
	}

	if err := viper.ReadInConfig(); err != nil {
		return newError(err, "Failed to read config file")
	}

	if err := viper.Unmarshal(&config); err != nil {
		return newError(err, "Failed to unmarshal config file")
	}

	return nil
}

func parseConfig(config *configuration) error {
	if config.TargetPath == "" {
		return newError(nil, "Target path has not been set")
	}

	if config.PkglistPath == "" {
		config.PkglistPath = config.TargetPath
	}

	for name, group := range config.Sources {
		if group.Path == "" {
			return newError(nil, fmt.Sprintf("%s has no defined path", name))
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

func verifyConfig(config *configuration) error {
	if err := checkDir(config.TargetPath); err != nil {
		return newError(err, "Failed to verify target path")
	}

	if err := checkDir(config.PkglistPath); err != nil {
		return newError(err, "Failed to verify package list path")
	}

	for name, group := range config.Sources {
		if err := checkDir(group.Path); err != nil {
			return newError(err, fmt.Sprintf("Failed to verify path for %s", name))
		}

		if err := checkPaths(group.Included, nil); err != nil {
			return newError(err, fmt.Sprintf("Failed to verify included path for %s", name))
		}

		if err := checkPaths(group.Excluded, nil); err != nil {
			return newError(err, fmt.Sprintf("Failed to verify excluded path for %s", name))
		}
	}

	return nil
}
