package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SyncFiles(config *configuration) error {
	dfs := config.allDotfiles

	for _, df := range dfs {
		if df.isUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", df.targetFile.name, df.lineCountDiff())

		if config.dryRun {
			continue
		}

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

func ShowDiffs(config *configuration) {
	dfs := config.allDotfiles

	for _, df := range dfs {
		if df.isUpToDate() {
			continue
		}

		if config.summaryView {
			fmt.Printf("%s (%+d lines)\n", df.targetFile.name, df.lineCountDiff())
		} else {
			fmt.Println(df.diff())
		}
	}
}
