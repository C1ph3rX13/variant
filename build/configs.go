package build

type CompileOpts struct {
	GoFileName  string // Go 文件名称
	ExeFileName string // Exe 文件名称
	CompilePath string // 编译路径
	HideConsole bool   // 编译隐藏控制台
	GDebug      bool   // 开启Debug日志
	GSeed       bool   // 随机 Base64 编码的种子
	Literals    bool   // 对字符串和数字字面量进行混淆
	Tiny        bool   // 最小化构建
	BuildMode   string // 构建模式, 推荐值:pie
}

type UpxOpts struct {
	Level   string // 压缩等级
	SrcExe  string // 目标文件
	SrcPath string // 压缩目录
	UpxPath string // upx.exe文件目录
	Keep    bool   // 保留原始文件
	Force   bool   // 强制压缩
}

type SignOpts struct {
	SignPath string // 签名路径
	UnSign   string // 未签名的原始文件
	Signed   string // 签名后的文件命名
	Thief    string // 窃取签名对象
	DstCert  string // 证书输出命名
	Cert     string // 指定证书
}
