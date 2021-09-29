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

func GetDiffs(config *configuration) []string {
	var diffs []string

	for _, df := range config.allDotfiles {
		if df.isUpToDate() {
			continue
		}

		if config.summaryView {
			diff := fmt.Sprintf("%s (%+d lines)", df.targetFile.name, df.lineCountDiff())
			diffs = append(diffs, diff)
		} else {
			diff := df.diff()
			diffs = append(diffs, diff)
		}
	}

	return diffs
}
