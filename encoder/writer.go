package encoder

import (
	"fmt"
	"os"
	"path/filepath"
	"variant/log"
)

func (p Payload) WriteStrings(content string) error {
	if content == "" {
		log.Fatal("payload is empty")
	}

	file, err := os.Create(filepath.Join(p.Path, p.FileName))
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
