package build

import (
	"fmt"
	"strings"
	"variant/rand"
)

// RenameGoTrimSuffix 删除文件名 .go 后缀，重命名为 .exe
func RenameGoTrimSuffix(name string) string {
	return fmt.Sprintf("%s.exe", strings.TrimSuffix(name, ".go"))
}

// RandomGoFile 随机命名 Go 代码文件
func RandomGoFile() string {
	return fmt.Sprintf("%s.go", rand.RStrings())
}

// RenameSignedPEName 命名签名后的PE文件
func RenameSignedPEName(name string) string {
	return fmt.Sprintf("signed_%s", name)
}
