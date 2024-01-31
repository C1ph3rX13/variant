package build

import (
	"fmt"
	"os"
	"variant/log"
)

func WriteStrings(content, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("<os.Create()> err: %s", err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("<WriteString()> err: %s", err)
	}

	return nil
}
