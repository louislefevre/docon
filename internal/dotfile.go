package internal

import (
	"bytes"
	"fmt"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
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
		if d.sourceFile.path == path || d.targetFile.path == path || d.targetFile.name == path {
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
	edits := myers.ComputeEdits(span.URIFromPath(df.targetFile.path),
		string(df.targetFile.contents), string(df.sourceFile.contents))
	diff := gotextdiff.ToUnified(df.targetFile.path, df.sourceFile.path, string(df.targetFile.contents), edits)
	return fmt.Sprint(diff)
}
