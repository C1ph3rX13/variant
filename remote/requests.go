package remote

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"variant/log"
)

// RestyGet Get 请求远程读取，返回 []byte 类型
func RestyGet(target string, proxyURL string) ([]byte, error) {
	c := CreateRestyClient()
	if proxyURL != "" {
		c.SetProxy(proxyURL)
	}

	resp, err := c.R().Get(target)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// RestyGetString 远程读取，返回 string 类型
func RestyGetString(target string, proxyURL string) string {
	body, err := RestyGet(target, proxyURL)
	if err != nil {
		log.Fatalf("request fail: %v", err)
	}

	return string(body)
}

// HttpGet 远程读取，返回 []byte 类型
func HttpGet(target string, proxyURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range UserAgentHeader() {
		req.Header.Set(key, value)
	}

	client := CreateHttpClient()
	if proxyURL != "" {
		proxyUrl, _ := url.Parse(proxyURL)
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

// HttpGetString 远程读取，返回 string 类型
func HttpGetString(target string, proxy string) string {
	body, err := HttpGet(target, proxy)
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
