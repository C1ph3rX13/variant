package sandbox

import (
	"golang.org/x/sys/windows/registry"
	"os"
)

func CheckWeChatExist() {
	// 打开注册表键 "SOFTWARE\Tencent\WeChat"
	reg, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Tencent\WeChat`, registry.QUERY_VALUE)
	if err != nil {
		os.Exit(0)
	}
	defer reg.Close()

	// 获取注册表键中的 "InstallPath" 值
	InstallPath, _, err := reg.GetStringValue("InstallPath")
	if err != nil || InstallPath == "" {
		os.Exit(0)
	}
}
