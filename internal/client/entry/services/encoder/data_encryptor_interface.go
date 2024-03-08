package encoder

type DataEncryptorInterface interface {
	Encrypt(data []byte, masterPass string) ([]byte, error)
	Decrypt(data []byte, masterPass string) ([]byte, error)
}
