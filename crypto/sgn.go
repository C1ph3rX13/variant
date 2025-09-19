package crypto

import (
	"fmt"

	sgn "github.com/EgeBalci/sgn/pkg"
)

// SgnEncoder 使用SGN对二进制文件进行编码
//
// 参数说明：
//   - file: 原始二进制数据
//   - arch: 种子值
//
// 返回值：
//   - []byte: 编码后的二进制数据
//   - error: 错误信息
func SgnEncoder(file []byte, arch int) ([]byte, error) {
	// 创建新的编码器实例，设置种子值
	encoder, err := sgn.NewEncoder(arch)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	err = encoder.SetArchitecture(arch)
	if err != nil {
		return nil, fmt.Errorf("failed to set architecture: %w", err)
	}

	encodedBinary, err := encoder.Encode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to encode binary: %w", err)
	}

	return encodedBinary, nil
}
