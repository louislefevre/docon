package internal

import (
	"errors"
	"fmt"
	"os"

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
