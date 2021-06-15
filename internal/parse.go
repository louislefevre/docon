package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func (dp dotfilePair) isUpToDate() bool {
	return bytes.Equal(dp.systemFile.contents, dp.targetFile.contents)
}

func (dp dotfilePair) lineCountDiff() int {
	systemFileCount := bytes.Count(dp.systemFile.contents, []byte{'\n'})
	repoFileCount := bytes.Count(dp.targetFile.contents, []byte{'\n'})
	return systemFileCount - repoFileCount
}

func parseConfiguration(config configuration) (dotfiles, error) {
	var dotfiles dotfiles

	for groupName, group := range config.mapping {
		if fileInfo, err := os.Stat(group.Path); os.IsNotExist(err) {
			return nil, err
		} else if !fileInfo.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", group.Path)
		}

		for _, file := range group.Included {
			err := filepath.Walk(file, visitCheck())
			if err != nil {
				return nil, err
			}
		}

		for _, file := range group.Excluded {
			err := filepath.Walk(file, visitCheck())
			if err != nil {
				return nil, err
			}
		}

		var files []string
		err := filepath.Walk(group.Path, visit(&files, group.Included, group.Excluded))
		if err != nil {
			return nil, err
		}

		for _, systemFilePath := range files {
			systemFileName := strings.ReplaceAll(systemFilePath, group.Path, "")
			targetFileName := filepath.Join(groupName, systemFileName)
			targetFilePath := filepath.Join(config.repoPath, targetFileName)

			systemFileContents, err := ioutil.ReadFile(systemFilePath)
			if err != nil {
				return nil, err
			}

			var repoFileContents []byte
			if _, err := os.Stat(targetFilePath); !os.IsNotExist(err) {
				repoFileContents, err = ioutil.ReadFile(targetFilePath)
				if err != nil {
					return nil, err
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
					contents: repoFileContents,
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
				fmt.Printf("Warning: file '%s' is both excluded and included\n", path)
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
