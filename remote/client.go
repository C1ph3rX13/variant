package remote

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/imroc/req/v3"
)

var headers = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

func CreateRestyClient() *resty.Client {
	client := resty.New().
		SetHeaders(headers).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(10 * time.Second).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)).
		SetRetryCount(3)

	return client
}

func CreateHttpClient() *http.Client {
	// 忽略证书
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client.Timeout = time.Second * 10

	return client
}

func CreateReqClient() *req.Client {
	client := req.C().
		SetTLSFingerprintChrome().
		EnableInsecureSkipVerify().
		DisableDebugLog()

	return client
}
