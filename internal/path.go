package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func visit(files *[]string, included []string, excluded []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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