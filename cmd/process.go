package cmd

import (
	"fmt"
	"path/filepath"
	"os"
)

func processConfiguration(config Configuration) error {
	for _, group := range config.ConfigMap {
		var files []string
		
		err := filepath.Walk(group.Path, visit(&files, group.Included, group.Excluded))
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Println(file)
		}
	}
	return nil
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
		
		if len(included) != 0 && !contains(included, path){
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
