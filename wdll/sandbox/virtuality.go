package sandbox

import (
	"os"
	"os/exec"
	"strings"
)

func WMICCheckVirtual() (bool, error) {
	// 执行命令获取计算机系统模型信息
	cmd := exec.Command("wmic", "path", "Win32_ComputerSystem", "get", "Model")
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// 将模型信息转换为小写字母方便匹配
	model := strings.ToLower(string(stdout))

	// 检查模型信息中是否包含虚拟机的关键词
	virtualKeywords := []string{"virtualbox", "virtual", "vmware", "kvm", "bochs", "hvm domu", "parallels"}
	for _, keyword := range virtualKeywords {
		if strings.Contains(model, keyword) {
			return true, nil
		}
	}

	return false, nil
}

// CheckVirtualFiles 调用 PathExists 检查沙箱或者虚拟机关键文件是否存在，如果存在则退出当前进程
func CheckVirtualFiles() {
	files := []string{
		"C:\\windows\\System32\\Drivers\\Vmmouse.sys",
		"C:\\windows\\System32\\Drivers\\vmtray.dll",
		"C:\\windows\\System32\\Drivers\\VMToolsHook.dll",
		"C:\\windows\\System32\\Drivers\\vmmousever.dll",
		"C:\\windows\\System32\\Drivers\\vmhgfs.dll",
		"C:\\windows\\System32\\Drivers\\vmGuestLib.dll",
		"C:\\windows\\System32\\Drivers\\VBoxMouse.sys",
		"C:\\windows\\System32\\Drivers\\VBoxGuest.sys",
		"C:\\windows\\System32\\Drivers\\VBoxSF.sys",
		"C:\\windows\\System32\\Drivers\\VBoxVideo.sys",
		"C:\\windows\\System32\\vboxdisp.dll",
		"C:\\windows\\System32\\vboxhook.dll",
		"C:\\windows\\System32\\vboxoglerrorspu.dll",
		"C:\\windows\\System32\\vboxoglpassthroughspu.dll",
		"C:\\windows\\System32\\vboxservice.exe",
		"C:\\windows\\System32\\vboxtray.exe",
		"C:\\windows\\System32\\VBoxControl.exe",
	}

	for _, file := range files {
		exists, _ := PathExists(file)
		if exists {
			os.Exit(0)
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
