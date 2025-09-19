package inject

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
	"variant/log"

	"golang.org/x/sys/windows"
)

func Reflective(dllPath string) {
	dllBytes, err := os.ReadFile(dllPath)
	if err != nil {
		log.Fatalf("Failed to open file %v", err)
	}
	dllPtr := uintptr(unsafe.Pointer(&dllBytes[0]))

	log.Infof("[+] DLL of size %d is loaded in memory at 0x%x\n", len(dllBytes), dllPtr)

	eLfanew := *((*uint32)(unsafe.Pointer(dllPtr + 0x3c)))
	ntHeader := (*IMAGE_NT_HEADERS64)(unsafe.Pointer(dllPtr + uintptr(eLfanew)))

	dllBase, err := windows.VirtualAlloc(uintptr(ntHeader.OptionalHeader.ImageBase),
		uintptr(ntHeader.OptionalHeader.SizeOfImage),
		windows.MEM_RESERVE|windows.MEM_COMMIT,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if err != nil {
		log.Fatalf("[!] VirtualAlloc Failed")
	}

	fmt.Printf("[+] Allocated address at 0x%x\n\n", dllBase)
	deltaImageBase := dllBase - uintptr(ntHeader.OptionalHeader.ImageBase)
	var numberOfBytesWritten uintptr
	err = windows.WriteProcessMemory(windows.CurrentProcess(), dllBase, &dllBytes[0], uintptr(ntHeader.OptionalHeader.SizeOfHeaders), &numberOfBytesWritten)
	if err != nil {
		log.Fatalf("[!] WriteProcessMemory Failed")
	}
	numberOfSections := int(ntHeader.FileHeader.NumberOfSections)

	var sectionAddr uintptr
	sectionAddr = dllPtr + uintptr(eLfanew) + unsafe.Sizeof(ntHeader.Signature) + unsafe.Sizeof(ntHeader.OptionalHeader) + unsafe.Sizeof(ntHeader.FileHeader)

	for i := 0; i < numberOfSections; i++ {
		section := (*IMAGE_SECTION_HEADER)(unsafe.Pointer(sectionAddr))
		sectionDestination := dllBase + uintptr(section.VirtualAddress)
		sectionBytes := (*byte)(unsafe.Pointer(dllPtr + uintptr(section.PointerToRawData)))
		fmt.Printf("[+] Copying %d bytes from 0x%x -> 0x%x for section : %s", section.SizeOfRawData, dllPtr+uintptr(section.PointerToRawData), sectionDestination, windows.ByteSliceToString(section.Name[:]))

		err = windows.WriteProcessMemory(windows.CurrentProcess(), sectionDestination, sectionBytes, uintptr(section.SizeOfRawData), &numberOfBytesWritten)
		if err != nil {
			log.Fatalf("[!] WriteProcessMemory Failed: %v \n", err)
		}

		fmt.Printf("	... Bytes 0x%x/0x%x Written\n", section.SizeOfRawData, numberOfBytesWritten)
		if windows.ByteSliceToString(section.Name[:]) == ".text" {
			var oldprotect uint32
			err := windows.VirtualProtect(sectionDestination, uintptr(section.SizeOfRawData), windows.PAGE_EXECUTE_READ, &oldprotect)
			if err != nil {
				log.Fatal("[ERROR] Failed to change memory permissions")
			}
		}
		sectionAddr += unsafe.Sizeof(*section)
	}

	relocations := ntHeader.OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_BASERELOC]
	relocationTable := uintptr(relocations.VirtualAddress) + dllBase
	fmt.Printf("[+] Relocation table Address: 0x%x\n\n", relocationTable)

	var relocationsProcessed int = 0
	for {

		relocationBlock := *(*BASE_RELOCATION_BLOCK)(unsafe.Pointer(uintptr(relocationTable + uintptr(relocationsProcessed))))
		relocEntry := relocationTable + uintptr(relocationsProcessed) + 8
		if relocationBlock.BlockSize == 0 && relocationBlock.PageAddress == 0 {
			break
		}
		relocationsCount := (relocationBlock.BlockSize - 8) / 2
		fmt.Printf("[+] PAGERVA : 0x%04x   Size: 0x%02x Entries Count: 0x%02x\n", relocationBlock.PageAddress, relocationBlock.BlockSize, relocationsCount)

		relocationEntries := make([]BASE_RELOCATION_ENTRY, relocationsCount)

		for i := 0; i < int(relocationsCount); i++ {
			relocationEntries[i] = *(*BASE_RELOCATION_ENTRY)(unsafe.Pointer(relocEntry + uintptr(i*2)))
		}
		for _, relocationEntry := range relocationEntries {
			if relocationEntry.Type() == 0 {
				continue
			}
			fmt.Printf("	--> Value: %X	Offset: %x\n", relocationEntry.OffsetType, relocationEntry.Offset())
			var size uintptr
			byteSlice := make([]byte, unsafe.Sizeof(size))
			relocationRVA := relocationBlock.PageAddress + uint32(relocationEntry.Offset())

			err = windows.ReadProcessMemory(windows.CurrentProcess(), dllBase+uintptr(relocationRVA), &byteSlice[0], unsafe.Sizeof(size), nil)
			if err != nil {
				log.Fatalf("[ERROR] Failed to ReadProcessMemory")
			}
			addressToPatch := uintptr(binary.LittleEndian.Uint64(byteSlice))
			addressToPatch += deltaImageBase
			a2Patch := uintptrToBytes(addressToPatch)
			err = windows.WriteProcessMemory(windows.CurrentProcess(), dllBase+uintptr(relocationRVA), &a2Patch[0], uintptr(len(a2Patch)), nil)
			if err != nil {
				log.Fatalf("[ERROR] Failed to WriteProcessMemory")
			}

		}
		relocationsProcessed += int(relocationBlock.BlockSize)

	}
	//time.Sleep(10 * time.Second)

	importsDirectory := ntHeader.OptionalHeader.DataDirectory[IMAGE_DIRECTORY_ENTRY_IMPORT]
	importDescriptorAddr := dllBase + uintptr(importsDirectory.VirtualAddress)
	log.Infof("[+] Import Descripton address: 0x%x\n\n", importDescriptorAddr)

	for {
		importDescriptor := *(*IMAGE_IMPORT_DESCRIPTOR)(unsafe.Pointer(importDescriptorAddr))
		if importDescriptor.Name == 0 {
			break
		}
		libraryName := uintptr(importDescriptor.Name) + dllBase
		dllName := windows.BytePtrToString((*byte)(unsafe.Pointer(libraryName)))
		fmt.Printf("[+] Importing DLL : %s\n", dllName)
		hLibrary, err := windows.LoadLibrary(dllName)
		if err != nil {
			log.Fatal("[ERROR] LoadLibrary Failed")
		}
		addr := dllBase + uintptr(importDescriptor.FirstThunk)
		for {
			thunk := *(*uint16)(unsafe.Pointer(addr))
			if thunk == 0 {
				break
			}
			functionNameAddr := dllBase + uintptr(thunk+2)

			functionName := windows.BytePtrToString((*byte)(unsafe.Pointer(functionNameAddr)))
			proc, err := windows.GetProcAddress(hLibrary, functionName)
			if err != nil {
				log.Fatal("[ERROR] Failed to GetProcAddress")
			}
			fmt.Printf("	--> Importing Function %s -> Addr: 0x%x\n", functionName, proc)
			procBytes := uintptrToBytes(proc)
			// https://reverseengineering.stackexchange.com/questions/16870/import-table-vs-import-address-table
			var numberOfBytesWritten uintptr
			err = windows.WriteProcessMemory(windows.CurrentProcess(), addr, &procBytes[0], uintptr(len(procBytes)), &numberOfBytesWritten)
			if err != nil {
				log.Fatal("[ERROR] Failed to WriteProcessMemory")
			}
			addr += 0x8

		}
		importDescriptorAddr += 0x14
	}
	//fmt.Printf("BreakPoint %x", dllBase+0x1251)
	//time.Sleep(time.Second * 10)
	syscall.SyscallN(dllBase+uintptr(ntHeader.OptionalHeader.AddressOfEntryPoint), dllBase, DLL_PROCESS_ATTACH, 0)
	fmt.Println("[+] DLL function executed")
	err = windows.VirtualFree(dllBase, 0x0, windows.MEM_RELEASE)
	if err != nil {
		log.Fatal("[ERROR] Failed to Free Memory")
	}
	fmt.Printf("[+] Freed Memory at 0x%x\n", dllBase)
}
