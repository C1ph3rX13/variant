package crypto

import (
	"fmt"
	"strconv"
	"strings"
)

var roarKeyArr = []string{"嗷", "呜", "啊", "~"}

// RoarEncode 编码字符串
func RoarEncode(a string) string {
	c := strings.Split(a, "")
	var d strings.Builder

	for i := 0; i < len(c); i++ {
		b := fmt.Sprintf("%x", c[i][0])
		for len(b) < 4 {
			b = "0" + b
		}
		d.WriteString(b)
	}

	dStr := d.String()
	var b strings.Builder

	for i := 0; i < len(dStr); i++ {
		cInt := int(dStr[i]-'0') + i%16
		if cInt >= 16 {
			cInt -= 16
		}
		b.WriteString(roarKeyArr[cInt/4] + roarKeyArr[cInt%4])
	}

	return roarKeyArr[3] + roarKeyArr[1] + roarKeyArr[0] + b.String() + roarKeyArr[2]
}

// RoarDecode 解码字符串
func RoarDecode(a string) string {
	if len(a) < 4 {
		return ""
	}

	c := roarKeyArr
	c[0] = a[2:3]
	c[1] = a[1:2]
	c[2] = a[len(a)-1 : len(a)]
	c[3] = a[0:1]
	a = a[3 : len(a)-1]

	var d strings.Builder
	b := strings.Split(a, "")
	for e := 0; e <= len(a)-2; e += 2 {
		var g, h int
		for g = 0; g < 4 && b[e] != c[g]; g++ {
		}
		for h = 0; h < 4 && b[e+1] != c[h]; h++ {
		}
		g = 4*g + h - (e/2)%16
		if g < 0 {
			g += 16
		}
		d.WriteString(fmt.Sprintf("%x", g))
	}

	var result strings.Builder
	for i := 0; i <= len(d.String())-4; i += 4 {
		e := d.String()[i : i+4]
		if parsedInt, err := strconv.ParseInt(e, 16, 64); err == nil {
			result.WriteString(string(rune(parsedInt)))
		} else {
			fmt.Println("Error parsing int:", err)
		}
	}

	return result.String()
}
