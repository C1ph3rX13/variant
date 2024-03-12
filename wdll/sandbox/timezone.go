package sandbox

import (
	"time"
)

// IsBeijingTimezone 用于判断当前是否处于北京时区
func IsBeijingTimezone() bool {
	// 加载北京时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false
	}

	// 获取当前时间
	now := time.Now()

	// 判断当前时间是否为北京时间
	return now.In(loc).Format("2006-01-02 15:04:05") == now.Format("2006-01-02 15:04:05")
}
