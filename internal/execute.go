package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

func ExecuteSync() {
	config, err := initConfig()
	checkErr(err)

	err = syncFiles(config)
	checkErr(err)
}

func ExecuteDiff(filePaths []string) {
	config, err := initConfig()
	checkErr(err)

	showDiffs(config, filePaths)
}

func ExecutePkg() {
	config, err := initConfig()
	checkErr(err)

	err = genPackageList(config)
	checkErr(err)
}

func ExecuteVerify() {
	_, err := initConfig()
	fmt.Printf("Using config file %s\n", viper.ConfigFileUsed())
	if err != nil {
		fmt.Println("Config file syntax is invalid.")
		checkErr(err)
	} else {
		fmt.Println("Config file syntax is valid.")
	}
}

func ExecuteCommit() {
	config, err := initConfig()
	checkErr(err)

	err = commitAll(config)
	checkErr(err)
}
