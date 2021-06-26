package internal

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
	dotfiles := ExecuteConfig()
	err := genPackageList(dotfiles)
	checkErr(err)
}
