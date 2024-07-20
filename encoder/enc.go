package encoder

import (
	"errors"
	"fmt"
	"reflect"
)

func (payload Payload) SetKeyIV(sign interface{}) (string, error) {
	if len(payload.PlainText) == 0 {
		return "", errors.New("plaintext is empty")
	}

	binRaw, err := BinReader(payload.PlainText)
	if err != nil {
		return "", fmt.Errorf("failed to read binary data: %w", err)
	}

	if len(payload.Key) != 16 && len(payload.IV) != 16 {
		return "", fmt.Errorf("the length of key and iv should be greater equal to 16")
	}

	// 获取签名函数的值和类型
	// Never obfuscate the Message type.
	signValue := reflect.ValueOf(sign)
	signType := signValue.Type()

	switch signType.Kind() {
	case reflect.Func:
		// 获取函数参数数量
		numParams := signType.NumIn()
		if numParams < 2 || numParams > 3 {
			return "", errors.New("sign function should have 2 or 3 parameters")
		}

		// 创建参数值切片
		params := make([]reflect.Value, numParams)
		params[0] = reflect.ValueOf(binRaw)
		params[1] = reflect.ValueOf(payload.Key)

		if numParams == 3 {
			params[2] = reflect.ValueOf(payload.IV)
		}

		// 调用函数并获取结果
		results := signValue.Call(params)

		if len(results) < 2 {
			return "", errors.New("sign function does not return expected result")
		}

		// 提取加密后的密文和可能的错误值
		cipherText := results[0].String()
		errValue := results[1].Interface()

		// 判断加密是否成功
		if errValue != nil {
			return "", fmt.Errorf("encryption failed: %v", errValue)
		}

		return cipherText, nil
	default:
		return "", errors.New("sign is not a valid function")
	}
}

func (payload Payload) NoKeyIV(sign interface{}) (string, error) {
	if len(payload.PlainText) == 0 {
		return "", errors.New("plaintext is empty")
	}

	binRaw, err := BinReader(payload.PlainText)
	if err != nil {
		return "", fmt.Errorf("failed to read binary data: %w", err)
	}

	// 获取签名函数的值和类型
	signValue := reflect.ValueOf(sign)
	signType := signValue.Type()

	switch signType.Kind() {
	case reflect.Func:
		// 获取函数参数数量
		numParams := signType.NumIn()
		if numParams < 0 || numParams > 2 {
			return "", errors.New("sign function should have 1 or 2 parameters")
		}

		// 创建参数值切片
		params := make([]reflect.Value, numParams)
		params[0] = reflect.ValueOf(binRaw)

		// 调用函数并获取结果
		results := signValue.Call(params)

		// 提取加密后的密文和可能的错误值
		cipherText := results[0].String()
		if len(results) >= 2 {
			errValue := results[1].Interface()
			// 判断加密是否成功
			if errValue != nil {
				return "", fmt.Errorf("encryption failed: %v", errValue)
			}
		}

		return cipherText, nil
	default:
		return "", errors.New("sign is not a valid function")
	}
}

func (payload Payload) PokemonStrings(sign interface{}) ([]string, error) {
	if len(payload.PlainText) == 0 {
		return nil, errors.New("plaintext is empty")
	}

	binRaw, err := BinReader(payload.PlainText)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary data: %w", err)
	}

	// 获取签名函数的值和类型
	signValue := reflect.ValueOf(sign)
	signType := signValue.Type()

	switch signType.Kind() {
	case reflect.Func:
		// 获取函数参数数量
		numParams := signType.NumIn()
		if numParams < 0 || numParams > 2 {
			return nil, errors.New("sign function should have 1 or 2 parameters")
		}

		// 创建参数值切片
		params := make([]reflect.Value, numParams)
		params[0] = reflect.ValueOf(binRaw)

		// 调用函数并获取结果
		results := signValue.Call(params)

		var cipherText []string
		// 提取加密后的密文和可能的错误值
		cipherText = results[0].Interface().([]string)

		return cipherText, nil
	default:
		return nil, errors.New("sign is not a valid function")
	}
}
