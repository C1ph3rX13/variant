package encoder

type Payload struct {
	PlainText string // 原始 Payload
	FileName  string // 加密后的文件
	Path      string
	Key       []byte
	IV        []byte
}
