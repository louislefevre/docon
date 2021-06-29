package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, item := range slice1 {
		if !containsString(slice2, item) {
			diff = append(diff, item)
		}
	}
	return diff
}

func isDisjoint(s1 []string, s2 []string) bool {
	for _, item := range s1 {
		if containsString(s2, item) {
			return false
		}
	}
	return true
}

func checkErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

func newError(err error, msg string) error {
	errMsg := "ERROR: "

	if err == nil && msg == "" {
		panic("Error contents cannot be empty")
	}
	if err != nil {
		errMsg += fmt.Sprintf("%s", err)
	}
	if msg != "" {
		errMsg += fmt.Sprintf("\n%s", msg)
	}

	return errors.New(color.RedString(errMsg))
}

func newWarning(msg string) string {
	if msg == "" {
		panic("Warning contents cannot be empty")
	}

	warnMsg := fmt.Sprintf("WARNING: %s", msg)
	return color.YellowString(warnMsg)
}

func checkPath(path string, check func(fs.FileInfo) bool) error {
	if fileInfo, err := os.Stat(path); err == nil {
		if check != nil && !check(fileInfo) {
			return fmt.Errorf("%s is an invalid path", path)
		}
		return nil
	} else {
		return err
	}
}

func checkFile(path string) error {
	return checkPath(path, func(info fs.FileInfo) bool {
		return info.Mode().IsRegular()
	})
}

func checkDir(path string) error {
	return checkPath(path, func(info fs.FileInfo) bool {
		return info.Mode().IsDir()
	})
}

func checkPaths(files []string, check func(fs.FileInfo) bool) error {
	fn := func(path string, _ os.FileInfo, _ error) error {
		err := checkPath(path, check)
		return err
	}

	for _, file := range files {
		if err := filepath.Walk(file, fn); err != nil {
			return err
		}
	}

	return nil
}
