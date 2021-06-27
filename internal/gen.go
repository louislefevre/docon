package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func genPackageList(dotfiles dotfiles) error {
	path, err := os.Getwd()
	if err != nil {
		return newError(err, "Failed to determine current working directory")
	}

	path = filepath.Join(path, "pkglist.txt")
	var targetContents []byte
	if _, err := os.Stat(path); err == nil {
		targetContents, err = ioutil.ReadFile(path)
		if err != nil {
			return newError(err, "Failed to read package list")
		}
	}

	cmd := exec.Command("uname", "-n")
	distro, err := cmd.Output()
	if err != nil {
		return newError(err, "Failed to determine operating system")
	}

	switch strings.Trim(string(distro), "\n") {
	case "arch":
		cmd = exec.Command("pacman", "-Q")
	default:
		return newError(nil, "Unknown operating system")
	}

	systemContents, err := cmd.Output()
	if err != nil {
		return newError(err, "Failed to determine installed packages")
	}

	targetSlice := strings.Split(string(targetContents), "\n")
	systemSlice := strings.Split(string(systemContents), "\n")
	added := difference(systemSlice, targetSlice)
	removed := difference(targetSlice, systemSlice)

	defer fmt.Printf("%d Packages\n", len(systemSlice))
	if len(added) == 0 && len(removed) == 0 {
		fmt.Println("Package list is up to date")
		return nil
	}

	if len(added) != 0 {
		fmt.Printf("Added/Updated (%+d):\n", len(added))
		for _, item := range added {
			fmt.Printf("- %s\n", strings.Fields(item)[0])
		}
	}

	if len(removed) != 0 {
		fmt.Printf("Removed (-%d):\n", len(removed))
		for _, item := range removed {
			fmt.Printf("- %s\n", strings.Fields(item)[0])
		}
	}

	err = ioutil.WriteFile(path, systemContents, 0644)
	if err != nil {
		return newError(err, "Failed to write to package list")
	}

	return nil
}
