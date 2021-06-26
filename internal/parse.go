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
		if fileInfo, err := os.Stat(group.Path); os.IsNotExist(err) {
			return nil, newError(err, fmt.Sprintf("%s does not exist", group.Path))
		} else if !fileInfo.IsDir() {
			return nil, newError(err, fmt.Sprintf("%s is not a directory", group.Path))
		}

		for _, file := range group.Included {
			err := filepath.Walk(file, visitCheck())
			if err != nil {
				return nil, newError(err, fmt.Sprintf("Failed to walk file tree for %s", file))
			}
		}

		for _, file := range group.Excluded {
			err := filepath.Walk(file, visitCheck())
			if err != nil {
				return nil, newError(err, fmt.Sprintf("Failed to walk file tree for %s", file))
			}
		}

		var files []string
		err := filepath.Walk(group.Path, visit(&files, group.Included, group.Excluded))
		if err != nil {
			return nil, newError(err, fmt.Sprintf("Failed to walk file tree for %s", group.Path))
		}

		for _, systemFilePath := range files {
			systemFileName := strings.ReplaceAll(systemFilePath, group.Path, "")
			targetFileName := filepath.Join(groupName, systemFileName)
			targetFilePath := filepath.Join(config.TargetPath, targetFileName)

			systemFileContents, err := ioutil.ReadFile(systemFilePath)
			if err != nil {
				return nil, newError(err, fmt.Sprintf("Failed to read file %s", systemFilePath))
			}

			var targetFileContents []byte
			if _, err := os.Stat(targetFilePath); err == nil {
				targetFileContents, err = ioutil.ReadFile(targetFilePath)
				if err != nil {
					return nil, newError(err, fmt.Sprintf("Failed to read file %s", systemFilePath))
				}
			}

			dotfiles = append(dotfiles, dotfilePair{
				systemFile: dotfile{
					name:     systemFileName,
					path:     systemFilePath,
					contents: systemFileContents,
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

func visitCheck() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
		return nil
	}
}

func visit(files *[]string, included []string, excluded []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(path); os.IsNotExist(err) {
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
