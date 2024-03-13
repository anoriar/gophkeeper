package encoder

//go:generate mockgen -source=data_encryptor_interface.go -destination=mock_data_encryptor/mock_data_encryptor.go -package=mock_data_encryptor
type DataEncryptorInterface interface {
	Encrypt(data []byte, masterPass string) ([]byte, error)
	Decrypt(data []byte, masterPass string) ([]byte, error)
}
