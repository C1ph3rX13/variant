package payloads

import (
	"fmt"
	"os"
	"path/filepath"
)

type PayloadCtx struct {
	Raw        string // 原始 Payload
	Encryption string // 加密 Payload
	SavePath   string // 保存加密后 Payload 为文件
	PlainText  string // 加密后的 Payload 变量名
	CipherText string // 解密后的 Payload 变量名
}

// Save 将加密后的 Payload 写入文件
func (p *PayloadCtx) Save(name string) (err error) {
	path := filepath.Join(p.SavePath, name)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %q: %w", p.SavePath, err)
	}
	defer file.Close()

	if _, err = file.WriteString(p.Encryption); err != nil {
		return fmt.Errorf("write string to %q: %w", p.SavePath, err)
	}

	return nil
}
