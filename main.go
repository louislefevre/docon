package main

import (
	"os"

	"github.com/louislefevre/docon/cmd"
)

func main() {
	if exitCode := cmd.ExecuteRoot(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
