package secret

import (
	"os"
)

const (
	secretDirname          = "./.data/secret/"
	authTokenFilename      = ".token"
	masterPasswordFilename = ".pass"
)

type SecretRepository struct {
}

func NewSecretRepository() *SecretRepository {
	return &SecretRepository{}
}

func (s SecretRepository) SaveAuthToken(token string) error {
	err := s.mkdir()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(secretDirname+authTokenFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(token))
	if err != nil {
		return err
	}

	return nil
}

func (s SecretRepository) GetAuthToken() (string, error) {
	file, err := os.OpenFile(secretDirname+authTokenFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content []byte
	_, err = file.Read(content)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s SecretRepository) SaveMasterPassword(pass string) error {
	err := s.mkdir()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(secretDirname+masterPasswordFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(pass))
	if err != nil {
		return err
	}

	return nil
}

func (s SecretRepository) GetMasterPassword() (string, error) {
	file, err := os.OpenFile(secretDirname+masterPasswordFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content []byte
	_, err = file.Read(content)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s SecretRepository) mkdir() error {
	err := os.MkdirAll(secretDirname, 0755)
	if err != nil {
		return err
	}
	return nil
}
