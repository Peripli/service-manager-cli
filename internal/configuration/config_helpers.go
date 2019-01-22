package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	defaultConfigDirName  = ".sm"
	defaultConfigFileName = "config.json"
)

func ensureDirExists(path string) error {
	dirPath := filepath.Dir(path)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if mkderr := os.Mkdir(dirPath, 0700); mkderr != nil {
			return mkderr
		}
	}
	return nil
}

func defaultFilePath() (string, error) {
	homeDir, err := homeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, defaultConfigDirName, defaultConfigFileName), nil
}

func homeDir() (string, error) {
	var homeDir string

	if os.Getenv("SM_HOME") != "" {
		homeDir = os.Getenv("SM_HOME")

		if _, err := os.Stat(homeDir); os.IsNotExist(err) {
			return "", fmt.Errorf("error locating SM_HOME folder '%s'", homeDir)
		}
	} else {
		homeDir = userHomeDir()
	}

	return homeDir, nil
}

// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
var userHomeDir = func() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}
