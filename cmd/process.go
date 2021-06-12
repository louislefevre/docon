package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Dotfile struct {
	SystemFileName     string
	SystemFilePath     string
	SystemFileContents []byte
	RepoFileName       string
	RepoFilePath       string
	RepoFileContents   []byte
}

func (dotfile Dotfile) IsUpToDate() bool {
	return bytes.Equal(dotfile.SystemFileContents, dotfile.RepoFileContents)
}

func (dotfile Dotfile) LineCountDiff() int {
	systemFileCount := bytes.Count(dotfile.SystemFileContents, []byte{'\n'})
	repoFileCount := bytes.Count(dotfile.RepoFileContents, []byte{'\n'})
	return systemFileCount - repoFileCount
}

func processConfiguration(config Configuration) ([]Dotfile, error) {
	var dotfiles []Dotfile
	for groupName, group := range config.ConfigMap {
		var files []string

		err := filepath.Walk(group.Path, visit(&files, group.Included, group.Excluded))
		if err != nil {
			return nil, err
		}

		for _, filePath := range files {
			fileName := strings.ReplaceAll(filePath, group.Path, "")
			repoName := filepath.Join(groupName, fileName)
			repoPath := filepath.Join(config.RepoPath, repoName)

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

			dotfiles = append(dotfiles, Dotfile{
				SystemFileName:     fileName,
				SystemFilePath:     filePath,
				RepoFileName:       repoName,
				RepoFilePath:       repoPath,
				SystemFileContents: fileContents,
				RepoFileContents:   repoContents,
			})
		}
	}
	return dotfiles, nil
}

func visit(files *[]string, included []string, excluded []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if contains(excluded, path) {
			if contains(included, path) {
				fmt.Printf("Warning: file '%s' is both excluded and included\n", path)
			}
			return nil
		}

		if len(included) != 0 && !contains(included, path) {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func syncFiles(dotfiles []Dotfile) error {
	for _, file := range dotfiles {
		if file.IsUpToDate() {
			continue
		}
		fmt.Printf("Updating %s (%+d lines)\n", file.RepoFileName, file.LineCountDiff())

		if _, err := os.Stat(file.RepoFilePath); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(file.RepoFilePath), os.ModePerm)
		} else if err != nil {
			return err
		}

		err := ioutil.WriteFile(file.RepoFilePath, file.SystemFileContents, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
