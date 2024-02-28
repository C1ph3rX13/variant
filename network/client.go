package network

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/imroc/req/v3"
)

func CreateRestyClient() *resty.Client {
	client := resty.New().
		SetHeaders(GetRandomAgent()).
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
		SetCommonHeaders(GetRandomAgent()).
		SetTLSFingerprintChrome().
		EnableInsecureSkipVerify().
		DisableDebugLog()

	return client
}
