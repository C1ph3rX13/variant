package remote

type Transfer struct {
	Src   string // 上传的加密文件
	Path  string // 加密文件文件的路径
	Proxy string // curl代理配置
}

type UsersCloud struct {
	Src  string // 上传的加密文件
	Path string // 加密文件文件的路径
}

type FileIO struct {
	Src  string // 上传的加密文件
	Path string // 加密文件文件的路径
}
