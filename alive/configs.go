package alive

import (
	"golang.org/x/sys/windows/registry"
)

const (
	Run             = `Software\Microsoft\Windows\CurrentVersion\Run`
	RunOnce         = `Software\Microsoft\Windows\CurrentVersion\RunOnce`
	RunServices     = `Software\Microsoft\Windows\CurrentVersion\RunServices`
	RunServicesOnce = `Software\Microsoft\Windows\CurrentVersion\RunServicesOnce`
)

type Registry struct {
	User    registry.Key // 用户、服务类型
	KeyName string       // 注册表键名
	CmdArgs string       // 注册表值名
	KeyPath string       // 注册表路径
}

type Program struct {
}
