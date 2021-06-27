package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func syncFiles(dfs dotfiles) error {
	for _, df := range dfs {
		if df.isUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", df.targetFile.name, df.lineCountDiff())

		if _, err := os.Stat(df.targetFile.path); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(df.targetFile.path), os.ModePerm)
		} else if err != nil {
			return newError(err, fmt.Sprintf("Failed to retrieve %s", df.targetFile.path))
		}

		err := ioutil.WriteFile(df.targetFile.path, df.sourceFile.contents, 0644)
		if err != nil {
			return newError(err, fmt.Sprintf("Failed to write to %s", df.targetFile.path))
		}
	}
	return nil
}

func showDiffs(dfs dotfiles, paths []string) {
	if len(paths) != 0 {
		for _, path := range paths {
			if df, ok := dfs.get(path); ok {
				fmt.Println(df.diff())
			} else {
				fmt.Printf("Could not show diff for %s: file is not being tracked\n", path)
			}
		}
	} else {
		for _, df := range dfs {
			fmt.Println(df.diff())
		}
	}
}
