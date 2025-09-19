package payloads

const (
	NullChunksize = 5
	PAYLOADCHUNK  = 10
)

func ReduceEntropy(payload []byte) ([]byte, int) {
	loop := len(payload) / PAYLOADCHUNK
	remainder := len(payload) % NullChunksize
	newPayloadSize := (3 * len(payload)) / 2

	newPayload := make([]byte, newPayloadSize+remainder)

	nPcntr, oPcntr := 0, 0

	for z := 0; z < loop; z++ {
		for i := 0; i < PAYLOADCHUNK; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}

		for j := 0; j < NullChunksize; j++ {
			newPayload[nPcntr] = 0x00
			nPcntr++
		}
	}

	if remainder > 0 {
		for i := 0; i != remainder; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}
	}

	return newPayload, newPayloadSize
}

func ReverseEntropy(payload []byte, payloadSize int) ([]byte, int) {
	remainder := payloadSize % NullChunksize
	newPayloadSize := (payloadSize / 3) * 2
	loop := newPayloadSize / PAYLOADCHUNK

	newPayload := make([]byte, newPayloadSize+remainder)

	nPcntr, oPcntr := 0, 0

	for i := 0; i < loop; i++ {
		for j := 0; j < PAYLOADCHUNK; j++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}

		for z := 0; z < NullChunksize; z++ {
			oPcntr++ // ignoring 5 bytes
		}
	}

	if remainder > 0 {
		for i := 0; i != remainder; i++ {
			newPayload[nPcntr] = payload[oPcntr]
			nPcntr++
			oPcntr++
		}
	}

	return newPayload, newPayloadSize
}
