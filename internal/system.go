package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func syncFiles(dotfiles dotfiles) error {
	for _, file := range dotfiles {
		if file.isUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", file.targetFile.name, file.lineCountDiff())

		if _, err := os.Stat(file.targetFile.path); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(file.targetFile.path), os.ModePerm)
		} else if err != nil {
			return err
		}

		err := ioutil.WriteFile(file.targetFile.path, file.systemFile.contents, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
