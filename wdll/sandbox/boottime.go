package sandbox

import (
	"os"
	"time"
	"variant/wdll"
)

func BootTime() {
	GetTickCount := wdll.GetTickCount()
	startTime, _, _ := GetTickCount.Call()
	if startTime == 0 {
		os.Exit(0)
	}

	checkTime := time.Duration(startTime * 1000 * 1000)
	setTime := 30 * time.Minute
	if checkTime < setTime {
		os.Exit(0)
	}
}
