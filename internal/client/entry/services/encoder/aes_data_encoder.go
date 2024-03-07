package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AesDataEncoder struct {
}

func NewAesDataEncoder() *AesDataEncoder {
	return &AesDataEncoder{}
}

func (d *AesDataEncoder) Encode(data []byte, key []byte) ([]byte, error) {

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
