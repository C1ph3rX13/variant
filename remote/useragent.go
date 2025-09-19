package remote

import (
	"fmt"
	"math/rand"
	"strings"
)

var browserTemplates = []struct {
	Name     string
	Template string
}{
	{
		Name:     "Chrome",
		Template: "Mozilla/5.0 ({os}) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{version} Safari/537.36",
	},
	{
		Name:     "Firefox",
		Template: "Mozilla/5.0 ({os}; rv:{version}) Gecko/20100101 Firefox/{version}",
	},
	{
		Name:     "Safari",
		Template: "Mozilla/5.0 ({os}) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/{version} Safari/605.1.15",
	},
	{
		Name:     "Edge",
		Template: "Mozilla/5.0 ({os}) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{version} Safari/537.36 Edg/{version}",
	},
	{
		Name:     "Brave",
		Template: "Mozilla/5.0 ({os}) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/{version} Safari/537.36",
	},
	{
		Name:     "Opera",
		Template: "Mozilla/5.0 ({os}) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{version} Safari/537.36 OPR/{version}",
	},
}

var osPool = []string{
	"Windows NT 10.0; Win64; x64",
	"Macintosh; Intel Mac OS X 10_15_7",
	"X11; Linux x86_64",
	"iPhone; CPU iPhone OS 15_0 like Mac OS X",
	"Android 10; Mobile Safari/537.36",
}

func randomVersion(browser string) string {
	switch browser {
	case "Chrome":
		return fmt.Sprintf("%d.0.0.0", rand.Intn(20)+100)
	case "Firefox":
		return fmt.Sprintf("%d.0", rand.Intn(20)+90)
	case "Safari":
		return fmt.Sprintf("%d.%d", rand.Intn(10)+15, rand.Intn(10)+1)
	case "Edge":
		return fmt.Sprintf("%d.0.0.0", rand.Intn(20)+100)
	case "Brave":
		return fmt.Sprintf("%d.0.0.0", rand.Intn(20)+90)
	case "Opera":
		return fmt.Sprintf("%d.0.0.0", rand.Intn(20)+90)
	default:
		return "1.0"
	}
}

// RandomUserAgent 生成随机User-Agent字符串
func RandomUserAgent() string {
	browser := browserTemplates[rand.Intn(len(browserTemplates))]
	os := osPool[rand.Intn(len(osPool))]
	version := randomVersion(browser.Name)

	ua := strings.ReplaceAll(browser.Template, "{os}", os)
	ua = strings.ReplaceAll(ua, "{version}", version)

	return ua
}

// UserAgentHeader 返回包含User-Agent的HTTP头
func UserAgentHeader() map[string]string {
	return map[string]string{
		"User-Agent": RandomUserAgent(),
	}
}
