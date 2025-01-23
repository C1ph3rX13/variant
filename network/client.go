package network

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/imroc/req/v3"
)

// 通用配置参数
const (
	requestTimeout = 10 * time.Second
	maxRedirects   = 10
	retryCount     = 3
)

// CreateRestyClient 创建预配置的Resty客户端
func CreateRestyClient() *resty.Client {
	return resty.New().
		SetHeaders(GetRandomAgent()).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(requestTimeout).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(maxRedirects)).
		SetRetryCount(retryCount)
}

// CreateHttpClient 创建标准库HTTP客户端
func CreateHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(*http.Request, []*http.Request) error { // 简化匿名函数
			return http.ErrUseLastResponse
		},
		Timeout: requestTimeout,
	}
}

// CreateReqClient 创建Req客户端
func CreateReqClient() *req.Client {
	return req.C().
		SetCommonHeaders(GetRandomAgent()).
		SetTLSFingerprintChrome().
		EnableInsecureSkipVerify().
		DisableDebugLog().
		SetTimeout(requestTimeout) // 显式设置超时保持配置一致性
}
