package sandbox

import (
	"os"
	"path/filepath"
)

func GetDesktopFiles() {
	desktopPath, err := os.UserHomeDir()
	if err != nil {
		os.Exit(0)
	}

	desktopFiles, err := os.ReadDir(filepath.Join(desktopPath, "Desktop"))
	if err != nil {
		os.Exit(0)
	}

	fileCount := len(desktopFiles)

	if fileCount < 10 {
		os.Exit(0)
	}
}
