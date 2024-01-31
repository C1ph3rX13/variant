package build

type CompileOpts struct {
	GoFileName  string // Go 文件名称
	ExeFileName string // Exe 文件名称
	HideConsole bool   // 编译隐藏控制台
	CompilePath string // 编译路径
}

type UpxOpts struct {
	Level   string // 压缩等级
	SrcExe  string // 目标文件
	Keep    bool
	Force   bool
	UpxPath string
}

type SignOpts struct {
	SignPath string // 签名路径
	UnSign   string // 未签名的原始文件
	Signed   string // 签名后的文件命名
	Thief    string // 窃取签名对象
	DstCert  string // 证书输出命名
	Cert     string // 指定证书
}
