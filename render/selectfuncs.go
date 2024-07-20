package render

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const (
	CryptoPath = "crypto"
)

func randSelect(functions []string) (string, error) {
	if len(functions) == 0 {
		return "", errors.New("没有找到符合条件的函数")
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := random.Intn(len(functions))
	return functions[randomIndex], nil
}

// SelectDecrypt 随机选择解密函数
func SelectDecrypt() (string, error) {
	functions, err := GetExportedFuncsFromFolder(CryptoPath, 1)
	if err != nil {
		return "", err
	}

	return randSelect(functions)
}

// SelectEncrypt PokemonDemo
func SelectEncrypt() (string, error) {
	functions, err := GetExportedFuncsFromFolder(CryptoPath, 0)
	if err != nil {
		return "", err
	}

	return randSelect(functions)
}

// SwapFuncs 匹配加解密函数
func SwapFuncs(funcName string) string {
	if strings.Contains(funcName, "Decrypt") {
		return strings.Replace(funcName, "Decrypt", "Encrypt", 1)
	} else {
		return strings.Replace(funcName, "Encrypt", "Decrypt", 1)
	}
}
