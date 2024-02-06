package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"variant/log"
)

// Upload 获取上传 file.io 后的文件Url
func (fi FileIO) Upload() string {
	url := fmt.Sprintf("https://file.io/?title=%s", fi.Src)

	resp, err := CreateRestyClient().
		SetHeader("accept", "application/json").
		SetHeader("Content-Type", " multipart/form-data").
		R().
		SetFile(fi.Src, filepath.Join(fi.Path, fi.Src)).
		SetContentLength(true).
		Post(url)
	if err != nil {
		log.Fatalf("upload request fail: %v", err)
	}

	if resp.StatusCode() == 200 && strings.Contains(resp.String(), "true") {
		link := getLink(resp.Body())
		return link
	}

	return ""
}

func getLink(body []byte) string {
	byteReader := bytes.NewReader(body)
	decoder := json.NewDecoder(byteReader)

	var result map[string]interface{}
	err := decoder.Decode(&result)
	if err != nil {
		log.Fatalf("Error occurred during decoding. Error: %s", err.Error())
	}

	if link, ok := result["link"].(string); ok {
		return link
	}

	return ""
}

// FileIORead 读取 file.io 的加密文件
func FileIORead(url string) string {
	resp, err := CreateRestyClient().
		SetHeader("Referer", "https://www.file.io/").
		R().Get(url)
	if err != nil {
		log.Fatalf("read request fail: %v", err)
	}

	return resp.String()
}
