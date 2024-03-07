package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
)

type DataEncoder struct {
	secretRepository secret.SecretRepositoryInterface
}

func NewDataEncoder(secretRepository secret.SecretRepositoryInterface) *DataEncoder {
	return &DataEncoder{secretRepository: secretRepository}
}

func (d *DataEncoder) Encode(data []byte, key []byte) ([]byte, error) {

	salt := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	encoded := append(salt, ciphertext...)

	return encoded, nil
}
