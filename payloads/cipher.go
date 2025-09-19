package payloads

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// Encrypt 使用反射调用任意加密函数，自动处理参数和返回值
//   - sign：加密类型函数/方法
//   - key ：加密使用的 Key
//   - iv ：加密使用的 IV
func (p *PayloadCtx) Encrypt(sign any, key, iv []byte) (any, error) {
	if len(p.Raw) == 0 {
		return nil, errors.New("payload is empty")
	}

	binRaw, err := os.ReadFile(p.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary data: %w", err)
	}

	signVal := reflect.ValueOf(sign)
	if signVal.Kind() != reflect.Func {
		return nil, errors.New("signature must be callable function")
	}

	results, err := ReflectCall(signVal, binRaw, key, iv)
	if err != nil {
		return nil, err
	}

	// 处理返回值
	if len(results) < 2 {
		return nil, errors.New("sign function does not return expected result")
	}

	result := results[0].Interface()
	errValue := results[1].Interface()

	// 返回值类型判断
	switch result.(type) {
	case string:
		if errValue != nil {
			return nil, fmt.Errorf("encryption failed: %w", errValue.(error))
		}
		return result.(string), nil
	case []string:
		if errValue != nil {
			return nil, fmt.Errorf("encryption failed: %w", errValue.(error))
		}
		return result.([]string), nil
	default:
		return nil, errors.New("unsupported return type")
	}
}

// ReflectCall 处理函数调用和参数构建
func ReflectCall(signVal reflect.Value, binRaw []byte, key []byte, iv []byte) ([]reflect.Value, error) {

	signType := signVal.Type()
	paramCount := signType.NumIn()

	if paramCount < 1 {
		return nil, errors.New("function requires at least one parameter")
	}

	params := make([]reflect.Value, 0, paramCount)
	params = append(params, reflect.ValueOf(binRaw))

	if paramCount > 1 {
		params = append(params, reflect.ValueOf(key))
	}
	if paramCount > 2 {
		params = append(params, reflect.ValueOf(iv))
	}

	results := signVal.Call(params)

	// 检查返回值数量
	if len(results) < 1 {
		return nil, errors.New("function returns no result")
	}

	return results, nil
}
