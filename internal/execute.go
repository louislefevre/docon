package internal

func ExecuteSync() {
	config, err := initConfig()
	checkErr(err)

	dotfiles, err := processConfiguration(config)
	checkErr(err)

	err = syncFiles(dotfiles)
	checkErr(err)
}
