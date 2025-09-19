package remote

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"variant/log"
)

const (
	// AESKeyLength AES密钥所需字节长度(16字节 = 128位)
	AESKeyLength = 16
	// DESKeyLength DES密钥所需字节长度(8字节 = 64位)
	DESKeyLength = 8
	// SHA256HexLength SHA-256哈希值的16进制表示长度
	SHA256HexLength = 64
)

const (
	BaiduIcoUrl      = "https://www.baidu.com/favicon.ico"
	ThreatbookIcoUrl = "https://x.threatbook.com/public/asset/img/favicon.ico"
	QianxinIcoUrl    = "https://www.qianxin.com/favicon.ico"
	AhIcoUrl         = "https://www.dbappsecurity.com.cn/images/favicon.ico"
	QhIcoUrl         = "https://www.360.cn/favicon.ico"
	QmIcoUrl         = "https://www.venustech.com.cn/r/cms/www/default/images/favicon.ico"
	CtIcoUrl         = "https://www.chaitin.cn/favicon.ico"
)

var ICONUrl = []string{
	BaiduIcoUrl,
	ThreatbookIcoUrl,
	QianxinIcoUrl,
	AhIcoUrl,
	QhIcoUrl,
	QmIcoUrl,
	CtIcoUrl,
}

// RandomICONUrl 随机选择一个图标URL
// 返回: 随机选择的图标URL字符串
func RandomICONUrl() string {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ICONUrl))))
	return ICONUrl[index.Uint64()]
}

// SHA256Hash 计算从指定URL获取内容的SHA-256哈希值
// 参数:
//
//	url: 要获取内容的URL地址
//	proxy: 可选的代理服务器地址
//
// 返回值:
//
//	[]byte: 内容的SHA-256哈希值(16进制格式)
func SHA256Hash(url string, proxy string) []byte {
	// 通过网络请求获取内容
	data, err := RestyGet(url, proxy)
	if err != nil {
		log.Fatalf("Get url %s failed: %v", url, err)
		return nil
	}

	// 计算SHA-256哈希值
	hash := sha256.Sum256(data)
	// 转换为16进制字符串并返回
	return []byte(fmt.Sprintf("%x", hash))
}

// ExtractAESKey 从SHA-256哈希值中提取128位(16字节)AES密钥
// 参数:
//
//	hash: SHA-256哈希值(16进制格式)
//	start: 密钥起始索引(包含)
//	end: 密钥结束索引(不包含)
//
// 返回值:
//
//	[]byte: 提取的AES密钥
//	error: 如果范围无效返回错误
func ExtractAESKey(hash []byte, start, end int) ([]byte, error) {
	// 验证范围是否有效
	if start < 0 || end > SHA256HexLength || end-start != AESKeyLength*2 {
		return nil, fmt.Errorf("AES密钥范围无效: 需要32个16进制字符(16字节), 实际获取 start=%d, end=%d", start, end)
	}
	return hash[start:end], nil
}

// ExtractDESKey 从SHA-256哈希值中提取64位(8字节)DES密钥
// 参数:
//
//	hash: SHA-256哈希值(16进制格式)
//	start: 密钥起始索引(包含)
//	end: 密钥结束索引(不包含)
//
// 返回值:
//
//	[]byte: 提取的DES密钥
//	error: 如果范围无效返回错误
func ExtractDESKey(hash []byte, start, end int) ([]byte, error) {
	// 验证范围是否有效
	if start < 0 || end > SHA256HexLength || end-start != DESKeyLength*2 {
		return nil, fmt.Errorf("DES密钥范围无效: 需要16个16进制字符(8字节), 实际获取 start=%d, end=%d", start, end)
	}
	return hash[start:end], nil
}
