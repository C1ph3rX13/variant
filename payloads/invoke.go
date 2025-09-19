package payloads

import (
	"fmt"
	"reflect"
)

func ReflectInvoke(fn any, args ...any) (any, error) {
	// 获取参数的反射值
	val := reflect.ValueOf(fn)

	// 检查是否为函数类型
	if val.Kind() != reflect.Func {
		return nil, fmt.Errorf("参数必须是函数或方法，实际类型: %T", fn)
	}

	// 检查参数数量是否匹配
	if len(args) != val.Type().NumIn() {
		return nil, fmt.Errorf("参数数量不匹配，期望 %d 个，实际 %d 个", val.Type().NumIn(), len(args))
	}

	// 准备参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		// 检查参数类型是否匹配
		argType := val.Type().In(i)
		argValue := reflect.ValueOf(arg)
		if !argValue.Type().AssignableTo(argType) {
			// 尝试类型转换
			if !argValue.Type().ConvertibleTo(argType) {
				return nil, fmt.Errorf("第 %d 个参数类型不匹配，期望 %s，实际 %s", i+1, argType, argValue.Type())
			}
			argValue = argValue.Convert(argType)
		}
		in[i] = argValue
	}

	// 调用函数
	rets := val.Call(in)

	// 处理返回值
	if len(rets) == 0 {
		return nil, nil
	}

	// 检查最后一个返回值是否为error类型
	lastRet := rets[len(rets)-1]
	if lastRet.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		// 如果是error类型
		errValue := lastRet.Interface()
		if errValue != nil {
			// 如果error不为nil，则返回错误
			return nil, errValue.(error)
		}
		// 如果error为nil，则去掉最后一个返回值
		rets = rets[:len(rets)-1]
	}

	// 将剩余返回值转换为[]any
	results := make([]any, len(rets))
	for i, ret := range rets {
		results[i] = ret.Interface()
	}

	return results, nil
}
