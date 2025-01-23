package persistence

import (
	"golang.org/x/sys/windows/registry"
)

// QueryRegistryValue 查询注册表键值
func (r Registry) QueryRegistryValue() (string, error) {
	key, err := registry.OpenKey(r.User, r.KeyPath, registry.ALL_ACCESS)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(r.KeyName)
	if err != nil {
		return "", err
	}

	return value, nil
}

// CreateRegistryKey 新增注册表键值
func (r Registry) CreateRegistryKey() error {
	key, _, err := registry.CreateKey(r.User, r.KeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue(r.KeyName, r.CmdArgs)
	if err != nil {
		return err
	}

	return nil
}
