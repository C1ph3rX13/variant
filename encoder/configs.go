package encoder

type Payload struct {
	PlainText string // 原始 Payload
	FileName  string // 加密后的shellcode文件名
	Path      string // 加密文件的保存路径
	Key       []byte
	IV        []byte
}
