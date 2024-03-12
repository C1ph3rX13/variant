package sandbox

import (
	"golang.org/x/sys/windows"
	"os"
)

func CheckLanguage() {
	systemLanguage, _ := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if len(systemLanguage) == 0 || systemLanguage[0] != "zh-CN" {
		os.Exit(0)
	}
}
