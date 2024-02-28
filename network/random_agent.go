package network

import (
	"math/rand"
	"time"
)

func GetRandomAgent() map[string]string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	browsers := map[string]string{
		"Chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
		"Safari":  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
		"Edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Brave":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/95.0.4638.54 Safari/537.36",
		"Vivaldi": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0 (Edition beta)",
		"Opera":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.0.0 (Edition beta)",
		"Whale":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Whale/2.10.123.42 Chrome/89.0.4389.128 Safari/537.36",
		"CocCoc":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Coc Coc/100.0.220 Chrome/96.0.4664.45 Safari/537.36",
	}

	selectedBrowser := random.Intn(len(browsers))
	counter := 0
	for _, value := range browsers {
		if counter == selectedBrowser {
			return map[string]string{"User-Agent": value}
		}
		counter++
	}

	return map[string]string{}
}
