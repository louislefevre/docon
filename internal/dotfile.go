package internal

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type dotfiles []dotfilePair

type dotfilePair struct {
	sourceFile dotfile
	targetFile dotfile
}

type dotfile struct {
	name     string
	path     string
	contents []byte
}

func (dfs dotfiles) get(path string) (dotfilePair, bool) {
	for _, d := range dfs {
		if d.sourceFile.path == path {
			return d, true
		} else if d.targetFile.path == path {
			return d, true
		}
	}
	return dotfilePair{}, false
}

func (dp dotfilePair) isUpToDate() bool {
	return bytes.Equal(dp.sourceFile.contents, dp.targetFile.contents)
}

func (dp dotfilePair) lineCountDiff() int {
	sourceFileCount := bytes.Count(dp.sourceFile.contents, []byte{'\n'})
	targetFileCount := bytes.Count(dp.targetFile.contents, []byte{'\n'})
	return sourceFileCount - targetFileCount
}

func (dp dotfilePair) diff() string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(dp.sourceFile.contents), string(dp.targetFile.contents), false)
	return dmp.DiffPrettyText(diffs)
}
