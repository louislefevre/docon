package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/louislefevre/docon/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("test", "scripts"),
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string] func() int{
		"docon": cmd.ExecuteRoot,
	}))
}
