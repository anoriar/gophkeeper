package encoder

type DataEncoderInterface interface {
	Encode(data []byte, key []byte) ([]byte, error)
}
