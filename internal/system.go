package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func syncFiles(dotfiles dotfiles) error {
	for _, pair := range dotfiles {
		if pair.isUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", pair.targetFile.name, pair.lineCountDiff())

		if _, err := os.Stat(pair.targetFile.path); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(pair.targetFile.path), os.ModePerm)
		} else if err != nil {
			return newError(err, fmt.Sprintf("Failed to retrieve %s", pair.targetFile.path))
		}

		err := ioutil.WriteFile(pair.targetFile.path, pair.systemFile.contents, 0644)
		if err != nil {
			return newError(err, fmt.Sprintf("Failed to write to %s", pair.targetFile.path))
		}
	}
	return nil
}

func showDiffs(dotfiles dotfiles, filePaths []string) {
	if len(filePaths) != 0 {
		for _, path := range filePaths {
			if pair, ok := dotfiles.get(path); ok {
				fmt.Println(pair.diff())
			} else {
				fmt.Printf("Could not show diff for %s: file is not being tracked\n", path)
			}
		}
	} else {
		for _, pair := range dotfiles {
			fmt.Println(pair.diff())
		}
	}
}
