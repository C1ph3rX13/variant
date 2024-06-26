package dynamic

import (
	"math/rand"
	"time"
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

var icons = map[string]string{
	BaiduIcoUrl:      "https://www.baidu.com/favicon.ico",
	ThreatbookIcoUrl: "https://x.threatbook.com/public/asset/img/favicon.ico",
	QianxinIcoUrl:    "https://www.qianxin.com/favicon.ico",
	AhIcoUrl:         "https://www.dbappsecurity.com.cn/images/favicon.ico",
	QhIcoUrl:         "https://www.360.cn/favicon.ico",
	QmIcoUrl:         "https://www.venustech.com.cn/r/cms/www/default/images/favicon.ico",
	CtIcoUrl:         "https://www.chaitin.cn/favicon.ico",
}

// RandomICO 随机选择一个 ICO 地址
func RandomICO() string {
	// 创建一个以当前时间为种子的随机数生成器
	rand.NewSource(time.Now().UnixNano())

	// 将 map 的 key 存入 slice 中
	keys := make([]string, 0, len(icons))
	for key := range icons {
		keys = append(keys, key)
	}

	// 随机选择一个 key
	randomKey := keys[rand.Intn(len(keys))]

	// 返回对应的 ICO 地址
	return icons[randomKey]
}
