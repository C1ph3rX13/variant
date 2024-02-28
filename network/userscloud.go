package network

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"variant/log"
)

func (uc UsersCloud) UCUpload() string {
	url := "https://u9458.userscloud.com/cgi-bin/upload.cgi?upload_type=file&utype=anon"
	resp, err := CreateReqClient().R().
		SetFile(uc.Src, filepath.Join(uc.Path, uc.Src)).
		Post(url)
	if err != nil {
		log.Fatalf("upload request fail: %v", err)
	}

	if resp.StatusCode == 200 && strings.Contains(resp.String(), "OK") {
		fileCode := getFileCode(resp.Bytes())
		log.Infof("UsersCloud file code: %s", fileCode)
		return fileCode
	}
	return ""
}

func getFileCode(body []byte) string {
	var result []map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("json.Unmarshal fail: %v", err)
	}

	if len(result) > 0 {
		return result[0]["file_code"].(string)
	}

	return ""
}

func UCRead(fn string) string {
	if fn == "" {
		log.Fatal("UsersCloud upload fail")
	}

	url := fmt.Sprintf("https://userscloud.com/%s", fn)
	body := fmt.Sprintf("op=download2&id=%s&rand=&referer=&method_free=&method_premium=&down_script=1&down_script=1", fn)

	p, err := CreateReqClient().R().
		SetBody(body).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Origin", "https://userscloud.com/").
		Post(url)
	if err != nil {
		os.Exit(0)
	}
	fmt.Println(p.Response.Status)
	fmt.Println(p.Response.Header)
	return p.String()
}
