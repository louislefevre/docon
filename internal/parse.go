package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func parseConfiguration(config configuration) (dotfiles, error) {
	var dotfiles dotfiles

	for groupName, group := range config.Sources {
		var files []string

		if err := filepath.Walk(group.Path, visit(&files, group.Included, group.Excluded)); err != nil {
			return nil, newError(err, fmt.Sprintf("Failed to walk file tree for %s", group.Path))
		}

		for _, path := range files {
			sourceFileName := strings.ReplaceAll(path, group.Path, "")
			sourceFilePath := path
			targetFileName := filepath.Join(groupName, sourceFileName)
			targetFilePath := filepath.Join(config.TargetPath, targetFileName)

			sourceFileContents, err := ioutil.ReadFile(sourceFilePath)
			if err != nil {
				return nil, newError(err, fmt.Sprintf("Failed to read file %s", sourceFilePath))
			}

			var targetFileContents []byte
			if err := checkFile(targetFilePath); err == nil {
				targetFileContents, err = ioutil.ReadFile(targetFilePath)
				if err != nil {
					return nil, newError(err, fmt.Sprintf("Failed to read file %s", targetFilePath))
				}
			}

			dotfiles = append(dotfiles, dotfilePair{
				sourceFile: dotfile{
					name:     sourceFileName,
					path:     sourceFilePath,
					contents: sourceFileContents,
				},
				targetFile: dotfile{
					name:     targetFileName,
					path:     targetFilePath,
					contents: targetFileContents,
				},
			})
		}
	}

	return dotfiles, nil
}

func visit(files *[]string, included []string, excluded []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err := checkPath(path, nil); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if containsString(excluded, path) {
			if containsString(included, path) {
				warning := newWarning(fmt.Sprintf("file '%s' is both excluded and included", path))
				fmt.Println(warning)
			}
			return nil
		}

		if len(included) != 0 && !containsString(included, path) {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}
