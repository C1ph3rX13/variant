package remote

import (
	"bytes"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"os"
)

func Resty(url string) ([]byte, error) {
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Http(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func SaveEncryptData(data []byte, fileName string) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
