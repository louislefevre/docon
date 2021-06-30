package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Checks if a string str is within slice s.
// Returns true if str is in s. Returns false otherwise.
func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Returns all elements in s1 that are not in s2.
func difference(s1 []string, s2 []string) []string {
	var diff []string
	for _, item := range s1 {
		if !containsString(s2, item) {
			diff = append(diff, item)
		}
	}
	return diff
}

// Checks if two slices share no elements.
// If the slices share any elements, returns false. Returns true otherwise.
func isDisjoint(s1 []string, s2 []string) bool {
	for _, item := range s1 {
		if containsString(s2, item) {
			return false
		}
	}
	return true
}

// Checks whether an error has occurred.
// Execution stops immediately if true. No-op if false.
func checkErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

// Returns a formatted error object.
// The err parameter can be nil, though execution will fail if msg is empty.
func newError(err error, msg string) error {
	if msg == "" {
		panic("Error contents cannot be empty")
	}

	errMsg := fmt.Sprintf("ERROR: %s", msg)
	if err != nil {
		errMsg += fmt.Sprintf("\n%s", err)
	}

	return errors.New(color.RedString(errMsg))
}

// Returns a formatted warning string.
// The err parameter can be nil, though execution will fail if msg is empty.
func newWarning(err error, msg string) string {
	if msg == "" {
		panic("Warning contents cannot be empty")
	}

	warnMsg := fmt.Sprintf("WARNING: %s", msg)
	if err != nil {
		warnMsg += fmt.Sprintf("\n%s", err)
	}

	return color.YellowString(warnMsg)
}

// Checks whether the provided path exists and is a file.
// Returns an error if it cannot be determined whether the path exists, or
// if the path is not a file.
func checkFile(path string) error {
	return checkPath(path, func(info fs.FileInfo) bool {
		return info.Mode().IsRegular()
	})
}

// Checks whether the provided path exists and is a directory.
// Returns an error if it cannot be determined whether the path exists, or
// if the path is not a directory.
func checkDir(path string) error {
	return checkPath(path, func(info fs.FileInfo) bool {
		return info.Mode().IsDir()
	})
}

// Checks whether the provided path exists.
// Returns an error if it cannot be determined whether the path exists, or
// if the check function returns false for a file.
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

// Checks whether each of the provided paths exists.
// Returns an error if it cannot be determined whether the path exists, or
// if the check function returns false for a file.
func checkPaths(paths []string, check func(fs.FileInfo) bool) error {
	for _, path := range paths {
		if err := checkPath(path, check); err != nil {
			return err
		}
	}
	return nil
}

// Checks each file in the file tree for the provided path.
// Returns an error if a file cannot be walked (e.g. if it doesn't exist), or
// if the check function returns false for a file.
func checkTree(path string, check func(fs.FileInfo) bool) error {
	fn := func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		err = checkPath(path, check)
		return err
	}
	if err := filepath.Walk(path, fn); err != nil {
		return err
	}
	return nil
}

// Checks each file in the file tree for each of the provided paths.
// Returns an error if a file cannot be walked (e.g. if it doesn't exist), or
// if the check function returns false for a file.
func checkTrees(paths []string, check func(fs.FileInfo) bool) error {
	for _, path := range paths {
		if err := checkTree(path, check); err != nil {
			return err
		}
	}
	return nil
}
