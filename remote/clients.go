package remote

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/imroc/req/v3"
)

// 通用安全配置参数
const (
	requestTimeout = 10 * time.Second // 请求超时时间
	maxRedirects   = 10               // 最大重定向次数
	retryCount     = 3                // 请求失败重试次数
)

// CreateRestyClient 创建预配置的Resty客户端（逻辑压缩版）
func CreateRestyClient() *resty.Client {
	return resty.New().
		SetHeaders(UserAgentHeader()).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(requestTimeout).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(maxRedirects)).
		SetRetryCount(retryCount)
}

// CreateHttpClient 创建标准库HTTP客户端（优化配置顺序）
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
// 采用链式调用优化，集成高级特性如TLS指纹模拟
func CreateReqClient() *req.Client {
	return req.C().
		// 设置随机User-Agent池
		SetCommonHeaders(UserAgentHeader()).
		// 模拟Chrome浏览器TLS指纹（防止反爬检测）
		SetTLSFingerprintChrome().
		// 跳过TLS验证
		EnableInsecureSkipVerify().
		// 禁用调试日志提升性能
		DisableDebugLog().
		// 统一设置请求超时时间
		SetTimeout(requestTimeout)
}
