package crypto

import (
	"encoding/hex"
	"net/url"
	"variant/log"
	"variant/remote"
)

const (
	Origin    = `http://hi.pcmoe.net`
	Referer   = `http://hi.pcmoe.net/buddha.html`
	BuddhaAPI = `http://hi.pcmoe.net/bear.php`
)

var BuddhaHeader = map[string]string{
	"Origin":           Origin,
	"Referer":          Referer,
	"X-Requested-With": "XMLHttpRequest",
	"X-Token":          "203B61D35068",
	"Content-type":     "application/x-www-form-urlencoded",
}

// sendRequest 发送请求并处理响应
func sendRequest(mode, code, txt string) (string, error) {
	v := url.Values{
		"mode": {mode},
		"code": {code},
		"txt":  {txt},
	}
	postBody := v.Encode()

	client := remote.CreateRestyClient()
	resp, respErr := client.R().
		SetHeaders(BuddhaHeader).
		SetBody(postBody).
		Post(BuddhaAPI)
	if respErr != nil {
		return "", respErr
	}

	return resp.String(), nil
}

// encode 编码明文
func encode(mode string, plainText []byte) string {
	hexText := hex.EncodeToString(plainText)
	encodedText, err := sendRequest(mode, "Encode", hexText)
	if err != nil {
		log.Fatal(err)
	}
	return encodedText
}

// decode 解码密文
func decode(mode string, cipherText string) []byte {
	decodedText, respErr := sendRequest(mode, "Decode", cipherText)
	if respErr != nil {
		log.Fatal(respErr)
	}
	stringText, _ := hex.DecodeString(decodedText)
	return stringText
}

// Deprecated: 加密接口废弃，等待更新
func BuddhaEncode(plainText []byte) string {
	return encode("Buddha", plainText)
}

// Deprecated: 加密接口废弃，等待更新
func BuddhaDecode(cipherText string) []byte {
	return decode("Buddha", cipherText)
}

// Deprecated: 加密接口废弃，等待更新
func BearEncode(plainText []byte) string {
	return encode("Bear", plainText)
}

// Deprecated: 加密接口废弃，等待更新
func BearDecode(cipherText string) []byte {
	return decode("Bear", cipherText)
}
