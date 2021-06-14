package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type dotfile struct {
	systemFileName     string
	systemFilePath     string
	systemFileContents []byte
	repoFileName       string
	repoFilePath       string
	repoFileContents   []byte
}

func (file dotfile) isUpToDate() bool {
	return bytes.Equal(file.systemFileContents, file.repoFileContents)
}

func (file dotfile) lineCountDiff() int {
	systemFileCount := bytes.Count(file.systemFileContents, []byte{'\n'})
	repoFileCount := bytes.Count(file.repoFileContents, []byte{'\n'})
	return systemFileCount - repoFileCount
}

func parseConfiguration(config configuration) ([]dotfile, error) {
	var dotfiles []dotfile

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

		for _, filePath := range files {
			fileName := strings.ReplaceAll(filePath, group.Path, "")
			repoName := filepath.Join(groupName, fileName)
			repoPath := filepath.Join(config.repoPath, repoName)

			fileContents, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			var repoContents []byte
			if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
				repoContents, err = ioutil.ReadFile(repoPath)
				if err != nil {
					return nil, err
				}
			}

			dotfiles = append(dotfiles, dotfile{
				systemFileName:     fileName,
				systemFilePath:     filePath,
				repoFileName:       repoName,
				repoFilePath:       repoPath,
				systemFileContents: fileContents,
				repoFileContents:   repoContents,
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
