package templates

import (
	"fmt"
	"reflect"
	"strings"
	"variant/compress"
	"variant/log"
	"variant/remote"
)

// Interval 生成数学区间表示字符串，格式为 "[start:end]"
//
// 注意:
//   - 本函数不对start/end的语义有效性进行检查
//   - 数值区间需自行保证字符串格式正确
func Interval(start int, end int) string {
	return fmt.Sprintf("[%v:%v]", start, end)
}

// AutoCompressPayload 使用指定的压缩算法对payload进行压缩，并返回压缩后的字符串。
// 参数：
//
//	payload    - 需要压缩的数据（字符串格式）
//	algorithm  - 压缩算法名称（如 "lzw", "zstd"），大小写不敏感
//	ratio      - 压缩参数（具体含义由算法决定，例如LZW压缩等级）
//
// 返回值：
//
//	压缩后的字符串，若压缩失败或算法不支持则返回空字符串
func AutoCompressPayload(payload string, algorithm string, ratio int) string {
	// 统一转换为小写以确保大小写不敏感的算法匹配
	algo := strings.ToLower(algorithm)

	var compressed string
	var err error

	// 根据算法类型执行对应压缩逻辑
	switch {
	case strings.Contains(algo, "lzw"):
		// 调用LZW算法压缩，使用传入的payload和ratio参数
		compressed, err = compress.LzwCompress(payload, ratio)
	case strings.Contains(algo, "zstd"):
		// 调用Zstandard算法压缩，当前实现不使用ratio参数
		compressed, err = compress.ZSTDCompress(payload)
	default:
		// 遇到不支持的算法时记录警告日志并直接返回空字符串
		log.Warnf("Unsupported compression algorithm: %q", algorithm)
		return ""
	}

	// 统一处理压缩过程中的错误
	if err != nil {
		log.Warnf("Compression failed with algorithm %q: %v", algorithm, err)
		return ""
	}

	return compressed
}

// DeriveKeyAndIVFromHash 根据指定的哈希算法对 iconUrl 进行哈希处理，
// 并基于 hash 值与给定的 keyRange 和 ivRange 拼接构造密钥 (key) 和初始化向量 (iv)。
//
// 参数：
//   - iconUrl: 要哈希处理的输入字符串（例如图标地址）
//   - hashFunc: 指定使用的哈希算法名称（如 "sha256" 或 "md5"），不区分大小写
//   - keyRange: key 的附加字段或偏移范围，用于构造最终密钥
//   - ivRange: iv 的附加字段或偏移范围，用于构造初始化向量
//
// 返回值：
//   - key: 构造出的密钥字符串
//   - iv: 构造出的初始化向量字符串
func DeriveKeyAndIVFromHash(iconUrl, hashFunc, keyRange, ivRange string) (key string, iv string) {
	lower := strings.ToLower(hashFunc)

	switch {
	case strings.Contains(lower, "sha256"):
		// 使用 SHA-256 算法生成哈希值
		hashValue := remote.SHA256Hash(iconUrl, "")

		// 将 hash 值与指定的 key/iv 范围拼接，生成最终的 key 和 iv
		key = fmt.Sprintf("%s%s", hashValue, keyRange)
		iv = fmt.Sprintf("%s%s", hashValue, ivRange)

		return key, iv

	case strings.Contains(lower, "md5"):
		// MD5 支持暂未实现，返回空值作为占位
		log.Warnf("MD5 hash function is not yet supported")
		return "", ""

	default:
		// 不支持的哈希算法类型，记录日志并返回空值
		log.Warnf("unsupported hash function: %q", hashFunc)
	}

	return "", ""
}

// SetZero  将传入的指针指向的值设置为其类型的零值
func SetZero[T any](v *T) {
	*v = zero[T]()
}

// 零值辅助函数
func zero[T any]() T {
	var z T
	return z
}

// FixSetZero 将任意类型的指针所指向的值设置为该类型的零值
func FixSetZero(v interface{}) error {
	val := reflect.ValueOf(v)

	// 确保输入是一个指针
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got %s", val.Kind())
	}

	// 获取指针指向的值
	elem := val.Elem()

	// 设置为该类型的零值
	elem.Set(reflect.Zero(elem.Type()))

	return nil
}
