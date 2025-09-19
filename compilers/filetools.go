package compilers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"variant/xrand"
)

// RenameGoTrimSuffix 删除文件名 .go 后缀，重命名为 .exe
func RenameGoTrimSuffix(file string) string {
	filename := filepath.Base(file)
	return fmt.Sprintf("%s.exe", strings.TrimSuffix(filename, ".go"))
}

// RandomGoFile 随机命名 Go 代码文件
func RandomGoFile() string {
	return fmt.Sprintf("%s.go", xrand.N(4))
}

// RenameSignedPEName 命名签名后的PE文件
func RenameSignedPEName(name string) string {
	return fmt.Sprintf("signed_%s", name)
}

// FindUpxBin 递归查找当前目录及其子目录中的 upx.exe 文件
func FindUpxBin() (string, error) {
	path, err := FindFile("Code.exe")
	if err != nil {
		return "", fmt.Errorf("FindFile failed: %v", err)
	}

	return path, err
}

// FindFile 递归查找当前目录及其子目录中指定的文件
//
// 参数:
//   - fileName: 要查找的文件名（如 "target.txt"）
//
// 返回:
//   - string: 找到的文件完整路径（格式：路径/文件名）
//   - error: 查找过程中的错误，未找到时返回 os.ErrNotExist
func FindFile(fileName string) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name is empty")
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	var foundPath string

	// 递归遍历文件系统
	err = filepath.WalkDir(wd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查是否匹配目标文件名
		if d.Name() == fileName {
			foundPath = path                            // 记录找到的路径
			return fmt.Errorf("file found at %s", path) // 用 error 终止遍历
		}

		return nil
	})

	if err != nil && foundPath == "" {
		if os.IsExist(err) {
			return "", fmt.Errorf("file %q not found", fileName)
		}
		return "", fmt.Errorf("search failed: %w", err)
	}

	if foundPath == "" {
		return "", fmt.Errorf("%w: file %q not found in %s", os.ErrNotExist, fileName, wd)
	}

	return foundPath, nil
}

// shouldSkipDir 判断是否跳过当前目录
func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "vendor", "node_modules", "__pycache__":
		return true
	default:
		return false
	}
}
