package sandbox

import (
	"os"
	"time"
	"variant/wdll"
)

/*
BeepSleep
使用 Beep 替代 Sleep

freq: 人类的听觉范围大约在 20Hz - 20000Hz 之间
duration : 设置时间，毫秒为单位
*/
func BeepSleep(duration uint32) {
	freq := uint32(30000)
	startTime := time.Now()

	r1, _, _ := wdll.Beep().Call(uintptr(freq), uintptr(duration))
	if r1 == 0 {
		os.Exit(0)
	}

	durationElapsed := time.Since(startTime).Seconds()

	if uint32(durationElapsed) != duration {
		os.Exit(0)
	}
}
