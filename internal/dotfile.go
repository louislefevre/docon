package internal

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type dotfiles []dotfile

type dotfile struct {
	sourceFile file
	targetFile file
}

type file struct {
	name     string
	path     string
	contents []byte
}

func (dfs dotfiles) get(path string) (dotfile, bool) {
	for _, d := range dfs {
		if d.sourceFile.path == path {
			return d, true
		} else if d.targetFile.path == path {
			return d, true
		}
	}
	return dotfile{}, false
}

func (df dotfile) isUpToDate() bool {
	return bytes.Equal(df.sourceFile.contents, df.targetFile.contents)
}

func (df dotfile) lineCountDiff() int {
	sourceFileCount := bytes.Count(df.sourceFile.contents, []byte{'\n'})
	targetFileCount := bytes.Count(df.targetFile.contents, []byte{'\n'})
	return sourceFileCount - targetFileCount
}

func (df dotfile) diff() string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(df.sourceFile.contents), string(df.targetFile.contents), false)
	return dmp.DiffPrettyText(diffs)
}
