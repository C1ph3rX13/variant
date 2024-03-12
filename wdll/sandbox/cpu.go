package sandbox

import (
	"os"
	"runtime"
)

func GetLogicalCPUCores() {
	CountOfCPU := runtime.NumCPU()
	if CountOfCPU < 4 {
		os.Exit(0)
	}
}
