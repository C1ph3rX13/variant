package encoder

import (
	"errors"
	"fmt"
	"reflect"
)

// SetKeyIV 根据传入的加密函数签名进行反射调用获取返回值
func (p Payload) SetKeyIV(sign interface{}) (string, error) {
	if len(p.PlainText) == 0 {
		return "", errors.New("plaintext is empty")
	}

	binRaw, err := BinReader(p.PlainText)
	if err != nil {
		return "", fmt.Errorf("failed to read binary data: %w", err)
	}

	// 获取签名函数的值和类型
	// 反射类型快速检查
	signVal := reflect.ValueOf(sign)
	if signVal.Kind() != reflect.Func {
		return "", errors.New("signature must be callable function")
	}

	// 获取函数参数数量
	signType := signVal.Type()
	if paramCount := signType.NumIn(); paramCount < 0 {
		return "", errors.New("function requires at least one parameter")
	}

	// 构建反射调用参数
	params := []reflect.Value{
		reflect.ValueOf(binRaw),
	}
	if signType.NumIn() > 1 {
		params = append(params, reflect.ValueOf(p.Key))
	}
	if signType.NumIn() > 2 {
		params = append(params, reflect.ValueOf(p.IV))
	}

	// 调用函数并获取结果
	results := signVal.Call(params)

	if len(results) < 2 {
		return "", errors.New("sign function does not return expected result")
	}

	// 类型转换
	cipherText, ok := results[0].Interface().(string)
	if !ok {
		return "", errors.New("invalid ciphertext type")
	}

	if errValue, errOk := results[1].Interface().(error); errOk && errValue != nil {
		return "", fmt.Errorf("encryption failed: %w", errValue)
	}

	return cipherText, nil
}

func (p Payload) PokemonEncoder(sign interface{}) ([]string, error) {
	if len(p.PlainText) == 0 {
		return nil, errors.New("plaintext is empty")
	}

	binRaw, err := BinReader(p.PlainText)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary data: %w", err)
	}

	signVal := reflect.ValueOf(sign)
	if signVal.Kind() != reflect.Func {
		return nil, errors.New("signature must be callable function")
	}

	// 获取函数参数数量
	signType := signVal.Type()
	if paramCount := signType.NumIn(); paramCount < 0 {
		return nil, errors.New("function requires at least one parameter")
	}

	// 构建反射调用参数
	params := []reflect.Value{
		reflect.ValueOf(binRaw),
	}

	// 调用函数并获取结果
	results := signVal.Call(params)

	// 提取加密后的密文和可能的错误值
	cipherText, ok := results[0].Interface().([]string)
	if !ok {
		return nil, errors.New("invalid Pokemon ciphertext type")
	}

	return cipherText, nil
}
