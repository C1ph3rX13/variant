package network

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"variant/log"
)

// Resty 远程读取，返回 []byte 类型
func Resty(url string, proxy string) ([]byte, error) {
	client := CreateRestyClient()
	if proxy != "" {
		client.SetProxy(proxy)
	}

	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	return body, nil
}

// RestyStrings 远程读取，返回 string 类型
func RestyStrings(url string, proxy string) string {
	body, err := Resty(url, proxy)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// Http 远程读取，返回 []byte 类型
func Http(target string, proxy string) ([]byte, error) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range GetRandomAgent() {
		req.Header.Set(key, value)
	}

	client := CreateHttpClient()
	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

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
func HttpStrings(target string, proxy string) string {
	body, err := Http(target, proxy)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// Req 远程读取，返回 []byte 类型
func Req(url string, proxy string) ([]byte, error) {
	client := CreateReqClient()
	if proxy != "" {
		client.SetProxyURL(proxy)
	}

	resp, err := client.R().
		SetRetryCount(5).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Bytes(), err
}

// ReqStrings 远程读取，返回 string 类型
func ReqStrings(url string, proxy string) string {
	body, err := Req(url, proxy)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// ReqIOReader 远程读取，返回 string 类型
func ReqIOReader(url string, proxy string) io.Reader {
	body, err := Req(url, proxy)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return bytes.NewReader(body)
}
