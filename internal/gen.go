package internal

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	pkgDateKeyword   keywordType = "{date}"
	pkgTimeKeyword   keywordType = "{time}"
	pkgSystemKeyword keywordType = "{system}"
)

var pkgKeywords = keywordSet{
	newKeyword(pkgDateKeyword),
	newKeyword(pkgTimeKeyword),
	newKeyword(pkgSystemKeyword),
}

func GenPackageList(config *configuration) error {
	var path string
	if config.Pkglist.Name != "" {
		path = filepath.Join(config.Pkglist.Path, config.Pkglist.Name)
	} else {
		path = filepath.Join(config.Pkglist.Path, "pkglist.txt")
	}

	for _, kw := range pkgKeywords {
		switch kw.kwType {
		case pkgDateKeyword:
			currentDate := time.Now().Local().Format("2006-01-02")
			path = kw.transform(path, currentDate)
		case pkgTimeKeyword:
			currentTime := time.Now().Local().Format("15:04:05")
			path = kw.transform(path, currentTime)
		case pkgSystemKeyword:
			if system, err := getUserOS(); err == nil {
				path = kw.transform(path, system)
			} else {
				return newError(err, "Failed to determine operating system")
			}
		}
	}

	var targetPackages []byte
	if err := checkFile(path); err == nil {
		targetPackages, err = ioutil.ReadFile(path)
		if err != nil {
			return newError(err, "Failed to read package list")
		}
	}

	systemPackages, err := getSystemPackages()
	if err != nil {
		return newError(err, "Failed to determine installed packages")
	}

	var (
		targetPackagesList  = strings.Split(string(targetPackages), "\n")
		systemPackagesList  = strings.Split(string(systemPackages), "\n")
		addedPackagesList   = difference(systemPackagesList, targetPackagesList)
		removedPackagesList = difference(targetPackagesList, systemPackagesList)
	)

	defer fmt.Printf("%d Packages\n", len(systemPackagesList))

	if len(addedPackagesList) != 0 {
		fmt.Printf("Added/Updated (%+d):\n", len(addedPackagesList))
		displayPackages(addedPackagesList)
	}

	if len(removedPackagesList) != 0 {
		fmt.Printf("Removed (-%d):\n", len(removedPackagesList))
		displayPackages(removedPackagesList)
	}

	if len(addedPackagesList) == 0 && len(removedPackagesList) == 0 {
		fmt.Println("Package list is up to date")
		return nil
	}

	if config.dryRun {
		return nil
	}

	if err := ioutil.WriteFile(path, systemPackages, 0644); err != nil {
		return newError(err, "Failed to write to package list")
	}

	return nil
}

func getSystemPackages() ([]byte, error) {
	userOS, err := getUserOS()
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	switch userOS {
	case "arch":
		cmd = exec.Command("pacman", "-Q")
	default:
		return nil, fmt.Errorf("unknown operating system: %s", userOS)
	}

	return cmd.Output()
}

func getUserOS() (string, error) {
	cmd := exec.Command("uname", "-n")
	if output, err := cmd.Output(); err == nil {
		return strings.Trim(string(output), "\n"), nil
	} else {
		return "", err
	}
}

func displayPackages(packages []string) {
	for _, item := range packages {
		fmt.Printf("- %s\n", strings.Fields(item)[0])
	}
}
