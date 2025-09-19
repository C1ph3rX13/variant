package sandbox

import (
	"os"
	"unsafe"
	"variant/xwindows"
)

func GetPhysicalMemory() {
	var memory uint64
	_, _ = xwindows.GetPhysicallyInstalledSystemMemory(uintptr(unsafe.Pointer(&memory)))

	memory = memory / 1048576
	if memory < 8 {
		os.Exit(0)
	}
}
