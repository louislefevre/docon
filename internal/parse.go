package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func gatherDotfiles(config *configuration) error {
	for name, group := range config.Sources {
		var dfs dotfiles
		var files []string

		if len(group.Included) != 0 {
			for _, path := range group.Included {
				if err := checkPath(path, nil); err != nil {
					return newError(err, fmt.Sprintf("Failed to find files for %s, missing path %s", name, path))
				}
				files = append(files, path)
			}
		} else {
			if err := filepath.Walk(group.Path, visit(&files, group.Excluded)); err != nil {
				return newError(err, fmt.Sprintf("Failed to walk file tree for %s", name))
			}
		}

		if name == "root" {
			name = ""
		}

		for _, path := range files {
			sourceFileName := strings.ReplaceAll(path, group.Path, "")
			sourceFilePath := path
			targetFileName := filepath.Join(name, sourceFileName)
			targetFilePath := filepath.Join(config.TargetPath, targetFileName)

			sourceFileContents, err := ioutil.ReadFile(sourceFilePath)
			if err != nil {
				return newError(err, fmt.Sprintf("Failed to read file %s", sourceFilePath))
			}

			var targetFileContents []byte
			if err := checkFile(targetFilePath); err == nil {
				targetFileContents, err = ioutil.ReadFile(targetFilePath)
				if err != nil {
					return newError(err, fmt.Sprintf("Failed to read file %s", targetFilePath))
				}
			}

			dfs = append(dfs, dotfile{
				sourceFile: file{
					name:     sourceFileName,
					path:     sourceFilePath,
					contents: sourceFileContents,
				},
				targetFile: file{
					name:     targetFileName,
					path:     targetFilePath,
					contents: targetFileContents,
				},
			})
		}
		group.dotfiles = dfs
		config.allDotfiles = append(config.allDotfiles, dfs...)
	}

	return nil
}

func visit(files *[]string, excluded []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if len(*files) > 100 {
			return fmt.Errorf("file tree is too large")
		} else if containsString(excluded, path) {
			return nil
		} else if info.IsDir() {
			return nil
		} else if err := checkPath(path, nil); err != nil {
			warning := newWarning(err, fmt.Sprintf("Failed to walk %s", path))
			fmt.Println(warning)
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}
