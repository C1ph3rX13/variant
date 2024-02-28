package network

import (
	"bytes"
	"io"
	"net/http"
	"variant/log"
)

// Resty 远程读取，返回 []byte 类型
func Resty(url string) ([]byte, error) {
	resp, err := CreateRestyClient().R().Get(url)
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
func RestyStrings(url string) string {
	body, err := Resty(url)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// Http 远程读取，返回 []byte 类型
func Http(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range GetRandomAgent() {
		req.Header.Set(key, value)
	}

	client := CreateHttpClient()
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

// HttpStrings 远程读取，返回 string 类型
func HttpStrings(url string) string {
	body, err := Http(url)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// Req 远程读取，返回 []byte 类型
func Req(url string) ([]byte, error) {
	resp, err := CreateReqClient().R().
		SetRetryCount(5).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Bytes(), err
}

// ReqStrings 远程读取，返回 string 类型
func ReqStrings(url string) string {
	body, err := Req(url)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}
