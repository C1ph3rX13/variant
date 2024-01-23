package enc

type Payload struct {
	PlainText string
	Key       []byte
	IV        []byte
}
