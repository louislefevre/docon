package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func syncFiles(dotfiles []dotfile) error {
	for _, file := range dotfiles {
		if file.isUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", file.repoFileName, file.lineCountDiff())

		if _, err := os.Stat(file.repoFilePath); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(file.repoFilePath), os.ModePerm)
		} else if err != nil {
			return err
		}

		err := ioutil.WriteFile(file.repoFilePath, file.systemFileContents, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}