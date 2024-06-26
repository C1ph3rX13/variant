package dynamic

import (
	"crypto/sha256"
	"fmt"
	"variant/log"
	"variant/network"
)

// CalcInternetResourcesHex 获取网络资源 SHA256 的指定切片
// 1字节(byte)等于8位(bits), 在16进制表示中，1个字节占用2个16进制位
func CalcInternetResourcesHex(url string, start, end int, proxies ...string) []byte {
	// 设置代理，仅第一个代理参数有效
	var proxy string
	if len(proxies) > 0 {
		proxy = proxies[0]
	} else {
		proxy = ""
	}

	data, err := network.Resty(url, proxy)
	if err != nil {
		log.Fatal("request fail: %v", err)
	}

	hash := sha256.Sum256(data)

	return []byte(fmt.Sprintf("%x", hash[start:end]))
}

func AesKey(url string, start, end int) []byte {
	if end-start != 8 || end > 64 {
		log.Fatal("invalid range, end - start = 8 || end <= 64")
	}

	hash := CalcInternetResourcesHex(url, start, end)

	return hash
}

func DesKey(url string, start, end int) []byte {
	if end-start != 4 || end > 64 {
		log.Fatal("invalid range, end - start = 4 || end <= 64")
	}

	hash := CalcInternetResourcesHex(url, start, end)

	return hash
}
