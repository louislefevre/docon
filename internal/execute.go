package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

func ExecuteSync(commitFiles bool) error {
	config, err := initConfig()
	if err != nil {
		return err
	}

	err = syncFiles(config)
	if err != nil {
		return err
	}

	if commitFiles {
		if err = commitAll(config); err != nil {
			return err
		}
	}

	return nil
}

func ExecuteDiff(filePaths []string) error {
	config, err := initConfig()
	if err != nil {
		return err
	}

	showDiffs(config, filePaths)
	return nil
}

func ExecutePkg() error {
	config, err := initConfig()
	if err != nil {
		return err
	}

	err = genPackageList(config)
	if err != nil {
		return err
	}

	return nil
}

func ExecuteVerify() error {
	_, err := initConfig()
	fmt.Printf("Using config file %s\n", viper.ConfigFileUsed())

	if err != nil {
		fmt.Println("Config file syntax is invalid.")
		return err
	} else {
		fmt.Println("Config file syntax is valid.")
		return nil
	}
}

func ExecuteCommit() error {
	config, err := initConfig()
	if err != nil {
		return err
	}

	err = commitAll(config)
	if err != nil {
		return err
	}

	return nil
}
