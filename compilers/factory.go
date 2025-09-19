package compilers

import (
	"fmt"
	"variant/gores"
	"variant/log"
)

// CompileCtx 编译配置
type CompileCtx struct {
	Build *BuildCtx       // 构建配置
	Upx   *UpxCtx         // UPX 压缩配置
	Sign  *CertThief      // 签名配置
	Res   *gores.GoWinRes // 资源配置
}

type BuildCtx struct {
	SrcGoFile   string // 模板渲染的 Go 代码文件
	OutExeFile  string // 编译后输出的 Exe 文件
	WorkingDir  string // 编译工作目录
	HideConsole bool   // 编译隐藏控制台
	BuildMode   string // 构建模式
	Obfuscate   bool   // Garble 编译
	Debug       bool   // 开启Debug日志
	Seed        bool   // 随机 Base64 编码的种子
	Literals    bool   // 对字符串和数字字面量进行混淆
	Tiny        bool   // 最小化构建
}

// NewGoBuildCtx 构建常规编译配置
//   - srcGoFile：模板渲染的 Go 代码文件
//   - outExeFile：编译后输出的 Exe 文件
//   - return: func(*BuildCtx)
func NewGoBuildCtx(srcGoFile string) func(*BuildCtx) {
	outExeFile := RenameGoTrimSuffix(srcGoFile)

	return func(b *BuildCtx) {
		b.SrcGoFile = srcGoFile
		b.OutExeFile = outExeFile
		b.BuildMode = "pie"
		b.WorkingDir = "output"
		b.HideConsole = false
	}
}

// NewGarbleBuildCtx 构建混淆编译配置
//   - srcGoFile：模板渲染的 Go 代码文件
//   - outExeFile：编译后输出的 Exe 文件
//   - return: func(*BuildCtx)
func NewGarbleBuildCtx(srcGoFile string) func(*BuildCtx) {
	outExeFile := RenameGoTrimSuffix(srcGoFile)

	return func(b *BuildCtx) {
		b.SrcGoFile = srcGoFile   // 模板渲染的 Go 代码文件
		b.OutExeFile = outExeFile // 编译后输出的 Exe 文件
		b.BuildMode = "pie"       // 编译模式
		b.WorkingDir = "output"   // 编译工作目录
		b.HideConsole = false     // 编译隐藏控制台
		b.Obfuscate = true        // 混淆编译，Garble
		b.Debug = false           // 开启Debug日志
		b.Seed = true             // 随机 Base64 编码的种子
		b.Literals = true         // 对字符串和数字字面量进行混淆
		b.Tiny = true             // 最小化构建
	}
}

type UpxCtx struct {
	Level      string // 压缩等级
	SrcExe     string // 目标文件
	WorkingDir string // 压缩文件路径
	UpxPath    string // upx.exe 文件目录
	Keep       bool   // 保留原始文件
	Force      bool   // 强制压缩
}

// NewUpxCtx 构建 Upx 压缩配置
//   - srcExe: 需要加壳的文件
//   - return: func(*UpxCtx)
//   - Note: 需要下载 upx.exe 放置于项目根目录下任意位置
func NewUpxCtx(srcExe string) func(*UpxCtx) {
	upxPath, err := FindUpxBin()
	if err != nil {
		log.Warnf("Failed to find upx binary: %v", err)
	}

	return func(u *UpxCtx) {
		u.Level = "--lzma"      // 默认压缩级别 --lzma
		u.SrcExe = srcExe       // 源文件名
		u.WorkingDir = "output" // 更标准的工作目录名称
		u.UpxPath = upxPath     // 设置找到的 upx 路径
		u.Keep = true           // 默认保留原始文件
		u.Force = false         // 默认不强制压缩
	}
}

// CertThief 结构体，用于签名操作
type CertThief struct {
	WorkingDir string // 签名目录
	SrcFile    string // 未签名的源文件
	Target     string // 已签名的目标文件
	SignedPE   string // 签名后的输出文件
	CertFile   string // 证书文件
}

// NewSignCtx 返回默认 Upx 压缩配置
//   - cert：指定的证书文件
//   - target：窃取签名的文件
//   - src：未签名的文件
//   - dst: 窃取签名后输出的文件
//   - return: func(*CertThief)
func NewSignCtx(target string, src string, dst string) func(*CertThief) {

	return func(s *CertThief) {
		s.SrcFile = src
		s.SignedPE = dst
		s.Target = target
		s.CertFile = "DstCert.cer"
		s.WorkingDir = "output"
	}
}

// NewResCtx 构建资源添加配置
//   - patchFile：需要添加资源文件的对象
//   - extractFile：需要提取资源文件的对象
//   - return: func(*gores.GoWinRes)
func NewResCtx(patchFile string, extractFile string) func(*gores.GoWinRes) {

	return func(g *gores.GoWinRes) {
		g.PatchDir = "output"
		g.PatchFile = patchFile
		g.ExtractDir = ""
		g.ExtractFile = extractFile
	}
}

type CompileBuilder struct {
	ctx *CompileCtx
}

func NewCompileBuilder() *CompileBuilder {
	return &CompileBuilder{
		ctx: &CompileCtx{},
	}
}

// SetBuildCtx 设置 BuildCtx
func (b *CompileBuilder) SetBuildCtx(fn func(*BuildCtx)) *CompileBuilder {
	if b.ctx.Build == nil {
		b.ctx.Build = &BuildCtx{}
	}

	fn(b.ctx.Build)
	return b
}

// SetUpxCtx 设置 UpxCtx
func (b *CompileBuilder) SetUpxCtx(fn func(*UpxCtx)) *CompileBuilder {
	if b.ctx.Upx == nil {
		b.ctx.Upx = &UpxCtx{}
	}

	fn(b.ctx.Upx)
	return b
}

// SetSignCtx 设置 SignCtx
func (b *CompileBuilder) SetSignCtx(fn func(*CertThief)) *CompileBuilder {
	if b.ctx.Sign == nil {
		b.ctx.Sign = &CertThief{}
	}

	fn(b.ctx.Sign)
	return b
}

// SetResCtx 设置 ResCtx
func (b *CompileBuilder) SetResCtx(fn func(res *gores.GoWinRes)) *CompileBuilder {
	if b.ctx.Res == nil {
		b.ctx.Res = &gores.GoWinRes{}
	}

	fn(b.ctx.Res)
	return b
}

// Build 返回最终构建的 CompileCtx（带字段验证）
func (b *CompileBuilder) Build() (*CompileCtx, error) {
	if b.ctx.Build == nil {
		return nil, fmt.Errorf("BuildCtx is not initialized")
	}

	if b.ctx.Build.SrcGoFile == "" {
		return nil, fmt.Errorf("SrcGoFile is required")
	}
	if b.ctx.Build.OutExeFile == "" {
		return nil, fmt.Errorf("OutExeFile is required")
	}
	if b.ctx.Build.WorkingDir == "" {
		return nil, fmt.Errorf("PatchDir is required")
	}

	return b.ctx, nil
}
