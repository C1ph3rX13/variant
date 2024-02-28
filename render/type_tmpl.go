package render

import (
	"fmt"
	"reflect"
	"text/template"
)

// RegisterTmplFuncs 注册自定义函数
func RegisterTmplFuncs() template.FuncMap {
	return template.FuncMap{
		"IsSlice":         IsSlice,
		"SliceType":       SliceType,
		"ConvertToString": ConvertToString,
	}
}

// IsSlice 检查变量是否为切片类型
func IsSlice(i interface{}) bool {
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Slice
}

// SliceType 获取切片元素的类型
func SliceType(i interface{}) []string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Slice {
		return []string{v.Type().Elem().String()}
	}
	return nil
}

// ConvertToString 将接口{}转换为字符串
func ConvertToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
