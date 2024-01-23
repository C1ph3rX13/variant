package render

import (
	"fmt"
	"os"
)

func WriteStringsToFile(data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("<os.Create()> err: %s", err)
	}

	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		return fmt.Errorf("<file.WriteString()> err: %s", err)
	}

	return nil
}

func ArgsValueToFile(data, filename string) error {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
