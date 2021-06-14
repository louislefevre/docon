package internal

func ExecuteSync() {
	config, err := initConfig()
	checkErr(err)

	dotfiles, err := parseConfiguration(config)
	checkErr(err)

	err = syncFiles(dotfiles)
	checkErr(err)
}
