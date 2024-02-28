package files

import (
	"os"
)

func Exists(src string) bool {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsFile(src string) (bool, error) {
	info, err := os.Stat(src)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return false, nil
	} else {
		return true, nil
	}
}

func IsDir(src string) (bool, error) {
	info, err := os.Stat(src)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func GetContent(src string) (string, error) {
	byteContent, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}

	return string(byteContent), nil
}

func WriteContent(filename string, text string) error {
	err := os.WriteFile(filename, []byte(text), 0644)

	if err != nil {
		return err
	}

	return nil
}

func Move(src string, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

func Copy(src string, dst string) error {
	check, err := IsFile(src)
	if err != nil {
		return err
	}

	if check == true {
		fileBytes, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		err = os.WriteFile(dst, fileBytes, 0644)
		if err != nil {
			return err
		}

	} else if check == false {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}

		err = os.MkdirAll(dst, srcInfo.Mode())
		if err != nil {
			return err
		}

		directory, _ := os.Open(src)
		objects, _ := directory.Readdir(-1)

		for _, obj := range objects {
			srcfilepointer := src + "/" + obj.Name()
			dstfilepointer := dst + "/" + obj.Name()

			if obj.IsDir() {
				err = Copy(srcfilepointer, dstfilepointer)
				if err != nil {
					return err
				}
			} else {
				fileBytes, err := os.ReadFile(srcfilepointer)
				if err != nil {
					return err
				}
				err = os.WriteFile(dstfilepointer, fileBytes, 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
