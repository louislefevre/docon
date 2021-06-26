package internal

import (
	"fmt"
	"os"
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
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
