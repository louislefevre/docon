package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lithammer/dedent"
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

// Removes intersecting values from two string slices.
// Values in s1 which are also in s2 are removed from s1.
// Returns s1 with the values removed.
func removeIntersecting(s1 []string, s2 []string) []string {
	var s3 []string
	for _, item := range s1 {
		if !containsString(s2, item) {
			s3 = append(s3, item)
		}
	}
	return s3
}

// Formats a multiline string by removing leading whitespace and linebreaks.
func multilineString(str string) string {
	str = strings.TrimPrefix(str, "\n")
	return dedent.Dedent(str)
}

// Reads data from standard input.
// Takes optional pretext strings as parameters, which will be printed prior
// to input being read.
func readStringInput(pretext ...string) string {
	for _, text := range pretext {
		fmt.Println(text)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input))
}

// Reads data from standard input and requests a boolean input until request is satisfied.
// Takes optional pretext strings as parameters, which will be printed prior
// to input being read.
func readBooleanInput(pretext ...string) bool {
	input := readStringInput(pretext...)

	switch input {
	case "yes", "y":
		return true
	case "no", "n":
		return false
	default:
		return readBooleanInput("Invalid input: enter yes/true or no/false.")
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

// Maps a slice of strings to corresponding groups, as identified by the @
// symbol (e.g. "bash@my_file"). The groups are returned as keys in a map,
// with their values as string slices containing each groups associated paths.
func splitGroupPaths(paths []string) (map[string][]string, error) {
	groupPaths := make(map[string][]string)

	for _, include := range paths {
		group := include
		path := ""

		if strings.Contains(include, "@") {
			if strings.Count(include, "@") > 1 {
				return nil, fmt.Errorf("can only specify one group per path")
			}
			split := strings.Split(include, "@")
			group, path = split[0], split[1]
		}

		if _, ok := groupPaths[group]; !ok {
			groupPaths[group] = []string{}
		}

		if path != "" {
			groupPaths[group] = append(groupPaths[group], path)
		}
	}

	return groupPaths, nil
}
