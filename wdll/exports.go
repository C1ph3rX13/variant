package wdll

const (
	// dll needs
	kernel32DLL = "kernel32.dll"
	ntdllDLL    = "ntdll.dll"
	rpcrt4DLL   = "Rpcrt4.dll"
	activedsDLL = "Activeds.dll"
	psapiDLL    = "psapi.dll"
	dbghelpDLL  = "dbghelp.dll"
	advapi32DLL = "advapi32.dll"

	// loader needs
	virtualAlloc                = "VirtualAlloc"
	virtualProtect              = "VirtualProtect"
	rtlCopyMemory               = "RtlCopyMemory"
	rtlCopyBytes                = "RtlCopyBytes"
	convertThreadToFiber        = "ConvertThreadToFiber"
	createFiber                 = "CreateFiber"
	switchToFiber               = "SwitchToFiber"
	getCurrentThread            = "GetCurrentThread"
	ntQueueApcThreadEx          = "NtQueueApcThreadEx"
	etwpCreateEtwThread         = "EtwpCreateEtwThread"
	waitForSingleObject         = "WaitForSingleObject"
	createThread                = "CreateThread"
	openProcess                 = "OpenProcess"
	virtualAllocEx              = "VirtualAllocEx"
	virtualProtectEx            = "VirtualProtectEx"
	writeProcessMemory          = "WriteProcessMemory"
	createRemoteThreadEx        = "CreateRemoteThreadEx"
	closeHandle                 = "CloseHandle"
	heapCreate                  = "HeapCreate"
	heapAlloc                   = "HeapAlloc"
	enumSystemLocalesA          = "EnumSystemLocalesA"
	uuidFromStringA             = "UuidFromStringA"
	getCurrentProcess           = "GetCurrentProcess"
	queueUserAPC                = "QueueUserAPC"
	allocADsMem                 = "AllocADsMem"
	enumSystemLocalesW          = "EnumSystemLocalesW"
	rtlEthernetAddressToStringA = "RtlEthernetAddressToStringA"
	rtlEthernetStringToAddressA = "RtlEthernetStringToAddressA"
	rtlIpv4StringToAddressA     = "RtlIpv4StringToAddressA"
	rtlIpv4AddressToStringA     = "RtlIpv4AddressToStringA"
	rtlMoveMemory               = "RtlMoveMemory"
	enumPageFilesW              = "EnumPageFilesW"
	enumerateLoadedModules      = "EnumerateLoadedModules"
	ntAllocateVirtualMemory     = "NtAllocateVirtualMemory"
	ntWriteVirtualMemory        = "NtWriteVirtualMemory"
	enumSystemLocalesEx         = "EnumSystemLocalesEx"
	reallocADsMem               = "ReallocADsMem"
	rtlCreateUserThread         = "RtlCreateUserThread"
	beep                        = "Beep"
	ntQueryInformationProcess   = "NtQueryInformationProcess"

	// etw
	etwEventWrite         = "EtwEventWrite"
	etwEventWriteEx       = "EtwEventWriteEx"
	etwEventWriteFull     = "EtwEventWriteFull"
	etwEventWriteString   = "EtwEventWriteString"
	etwEventWriteTransfer = "EtwEventWriteTransfer"

	// phant0m
	ntQueryInformationThread = "NtQueryInformationThread"
	i_QueryTagInformation    = "I_QueryTagInformation"
	openThread               = "OpenThread"
	terminateThread          = "TerminateThread"
	readProcessMemory        = "ReadProcessMemory"
	createToolhelp32Snapshot = "CreateToolhelp32Snapshot"
	thread32First            = "Thread32First"

	// sandbox needs
	getTickCount                       = "GetTickCount"
	getPhysicallyInstalledSystemMemory = "GetPhysicallyInstalledSystemMemory"
)

type IMAGE_DOS_HEADER struct { // DOS .EXE header
	/*E_magic    uint16     // Magic number
	  E_cblp     uint16     // Bytes on last page of file
	  E_cp       uint16     // Pages in file
	  E_crlc     uint16     // Relocations
	  E_cparhdr  uint16     // Size of header in paragraphs
	  E_minalloc uint16     // Minimum extra paragraphs needed
	  E_maxalloc uint16     // Maximum extra paragraphs needed
	  E_ss       uint16     // Initial (relative) SS value
	  E_sp       uint16     // Initial SP value
	  E_csum     uint16     // Checksum
	  E_ip       uint16     // Initial IP value
	  E_cs       uint16     // Initial (relative) CS value
	  E_lfarlc   uint16     // File address of relocation table
	  E_ovno     uint16     // Overlay number
	  E_res      [4]uint16  // Reserved words
	  E_oemid    uint16     // OEM identifier (for E_oeminfo)
	  E_oeminfo  uint16     // OEM information; E_oemid specific
	  E_res2     [10]uint16 // Reserved words*/
	E_lfanew uint32 // File address of new exe header
}

type IMAGE_NT_HEADER struct {
	Signature      uint32
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER
}

type IMAGE_FILE_HEADER struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type IMAGE_OPTIONAL_HEADER struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]IMAGE_DATA_DIRECTORY
}

type IMAGE_DATA_DIRECTORY struct {
	VirtualAddress uint32
	Size           uint32
}
