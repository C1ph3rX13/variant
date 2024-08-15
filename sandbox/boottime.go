package sandbox

import (
	"os"
	"time"
	"variant/xwindows"
)

func BootTime() {
	startTime, _ := xwindows.GetTickCount()
	if startTime == 0 {
		os.Exit(0)
	}

	checkTime := time.Duration(startTime * 1000 * 1000)
	setTime := 30 * time.Minute
	if checkTime < setTime {
		os.Exit(0)
	}
}

func BootTimeGetTime() {
	overtime, _ := xwindows.TimeGetTime()

	if overtime/3600000.0 < 9 {
		os.Exit(0)
	}
}
