package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

type AesDataEncryptor struct {
}

func NewAesDataEncoder() *AesDataEncryptor {
	return &AesDataEncryptor{}
}

func (d *AesDataEncryptor) Encrypt(data []byte, masterPass string) ([]byte, error) {

	key := d.createKeyFromMasterPass(masterPass)

	aes, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}
	encrypted := gcm.Seal(nonce, nonce, []byte(data), nil)

	return encrypted, nil
}

func (d *AesDataEncryptor) Decrypt(data []byte, masterPass string) ([]byte, error) {

	key := d.createKeyFromMasterPass(masterPass)

	aes, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	decrypted, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}
	return decrypted, nil
}

func (d *AesDataEncryptor) createKeyFromMasterPass(masterPass string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(masterPass))
	key := hasher.Sum(nil)
	return key
}
