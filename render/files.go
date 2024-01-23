package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetTmplFiles(tmplDir string) ([]string, error) {
	var tmplFiles []string

	files, err := os.ReadDir(tmplDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".render") {
			tmplFiles = append(tmplFiles, filepath.Join(tmplDir, file.Name()))
		}
	}

	return tmplFiles, nil
}

func GetBinFiles(binDir string) ([]string, error) {
	var binFiles []string

	files, err := os.ReadDir(binDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".bin") {
			binFiles = append(binFiles, filepath.Join(binDir, file.Name()))
		}
	}

	return binFiles, nil
}
