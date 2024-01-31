package remote

import (
	"bytes"
	"io"
	"net/http"
)

// Resty 远程读取，返回 []byte 类型
func Resty(url string) ([]byte, error) {
	resp, err := restyConf().R().Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RestyStrings 远程读取，返回 string 类型
func RestyStrings(url string) (string, error) {
	body, err := Resty(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Http 远程读取，返回 []byte 类型
func Http(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := httpConf()
	resp, err := client.Do(req)
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

// HttpString 远程读取，返回 string 类型
func HttpString(url string) (string, error) {
	body, err := Http(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
