package internal

import "fmt"

func ExecuteConfig() dotfiles {
	config, err := initConfig()
	checkErr(err)

	dotfiles, err := parseConfiguration(config)
	checkErr(err)

	return dotfiles
}

func ExecuteSync() {
	dotfiles := ExecuteConfig()
	err := syncFiles(dotfiles)
	checkErr(err)
}

func ExecuteDiff(filePaths []string) {
	dotfiles := ExecuteConfig()
	showDiffs(dotfiles, filePaths)
}

func ExecutePkg() {
	config, err := initConfig()
	checkErr(err)
	
	err = genPackageList(config)
	checkErr(err)
}

func ExecuteVerify() {
	_, err := initConfig()
	if err != nil {
		fmt.Println("Config file syntax is invalid.")
		checkErr(err)
	} else {
		fmt.Println("Config file syntax is valid.")
	}
}
