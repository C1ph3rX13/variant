package sandbox

import (
	"variant/xwindows"

	"github.com/gonutz/ide/w32"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// HideConsoleWin 已被沙箱标记为 CobaltStrike
func HideConsoleWin() {
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}

// HideConsoleW32 已被沙箱标记为 CobaltStrike
func HideConsoleW32() {
	hide := w32.GetConsoleWindow()
	if hide != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(hide)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(hide, w32.SW_HIDE)
		}
	}
}

// HideConsoleApi 已被沙箱标记为 CobaltStrike
func HideConsoleApi() {
	handle, _ := xwindows.GetConsoleWindow()

	err := xwindows.ShowWindow(windows.Handle(handle), 0)
	if err != nil {
		panic(err)
	}
}
