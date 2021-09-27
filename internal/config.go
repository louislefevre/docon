package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type configuration struct {
	TargetPath  string                   `mapstructure:"target"`
	PkglistPath string                   `mapstructure:"pkglist"`
	Git         gitStruct                `mapstructure:"git"`
	Sources     map[string]*sourceStruct `mapstructure:"sources"`
	allDotfiles dotfiles
}

type gitStruct struct {
	Dir       bool   `mapstructure:"dir"`
	CommitMsg string `mapstructure:"msg"`
	Author    struct {
		name  string `mapstructure:"name"`
		email string `mapstructure:"email"`
	}
}

type sourceStruct struct {
	Path      string   `mapstructure:"path"`
	CommitMsg string   `mapstructure:"msg"`
	Ignore    bool     `mapstructure:"ignore"`
	Included  []string `mapstructure:"include"`
	Excluded  []string `mapstructure:"exclude"`
	dotfiles  dotfiles
}

func InitConfig() (*configuration, error) {
	var config configuration

	if err := loadConfig(&config); err != nil {
		return &config, err
	}

	if err := applyFlags(&config); err != nil {
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

func createConfig() error {
	template := multilineString(`
	sources:
	  docon:
	    include:
	    - config.yaml
	    path: /home/.config/docon/
	target: /path/to/target
	`)

	if err := viper.ReadConfig(bytes.NewBuffer([]byte(template))); err != nil {
		return err
	}

	// TODO: Check that directory path exists before writing.
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	viper.WatchConfig()
	fmt.Printf("Created configuration file at %s\n", viper.ConfigFileUsed())
	fmt.Println("Modify the files contents to specify your configuration settings")
	return nil
}

func loadConfig(config *configuration) error {
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		if home, err := homedir.Dir(); err == nil {
			viper.AddConfigPath(fmt.Sprintf("%s/.config/docon/", home))
		} else {
			warning := newWarning(err, "Failed to find home directory")
			fmt.Println(warning)
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			dialog := "Config file not found - would you like to create one? [yes/no]"
			if ok := readBooleanInput(dialog); ok {
				if err := createConfig(); err != nil {
					return newError(err, "Failed to create config file")
				}
			}
			os.Exit(0)
		} else {
			return newError(err, "Failed to read config file")
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return newError(err, "Failed to unmarshal config file")
	}

	return nil
}

func applyFlags(config *configuration) error {
	if commitMsg := viper.GetString("message"); commitMsg != "" {
		config.Git.CommitMsg = commitMsg
	}

	if included := viper.GetStringSlice("only"); len(included) > 0 {
		for name, source := range config.Sources {
			if !containsString(included, name) {
				source.Ignore = true
			}
		}
	}

	if excluded := viper.GetStringSlice("ignore"); len(excluded) > 0 {
		for name, source := range config.Sources {
			if containsString(excluded, name) {
				source.Ignore = true
			}
		}
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

	if item, ok := containsValidKeywords(config.Git.CommitMsg, gitKeywords); !ok {
		return newError(nil, fmt.Sprintf("Git message has invalid keyword %s", item))
	}

	for name, group := range config.Sources {
		if group.Ignore {
			delete(config.Sources, name)
			continue
		}

		if strings.TrimSpace(name) == "" {
			return newError(nil, "Group name cannot be empty")
		}

		if group.Path == "" {
			return newError(nil, fmt.Sprintf("%s has no defined path", name))
		}

		if item, ok := containsValidKeywords(group.CommitMsg, gitKeywords); !ok {
			return newError(nil, fmt.Sprintf("Git message for %s has invalid keyword %s", name, item))
		}

		if !isDisjoint(group.Included, group.Excluded) {
			return newError(nil, fmt.Sprintf("%s contains items both included and excluded", name))
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
