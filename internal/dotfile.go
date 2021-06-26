package internal

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type dotfiles []dotfilePair

type dotfilePair struct {
	systemFile dotfile
	targetFile dotfile
}

type dotfile struct {
	name     string
	path     string
	contents []byte
}

func (dfs dotfiles) get(path string) (dotfilePair, bool) {
	for _, d := range dfs {
		if d.systemFile.path == path {
			return d, true
		} else if d.targetFile.path == path {
			return d, true
		}
	}
	return dotfilePair{}, false
}

func (dp dotfilePair) isUpToDate() bool {
	return bytes.Equal(dp.systemFile.contents, dp.targetFile.contents)
}

func (dp dotfilePair) lineCountDiff() int {
	systemFileCount := bytes.Count(dp.systemFile.contents, []byte{'\n'})
	repoFileCount := bytes.Count(dp.targetFile.contents, []byte{'\n'})
	return systemFileCount - repoFileCount
}

func (dp dotfilePair) diff() string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(dp.systemFile.contents), string(dp.targetFile.contents), false)
	return dmp.DiffPrettyText(diffs)
}
