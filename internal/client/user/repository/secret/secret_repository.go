package secret

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrTokenNotFound = errors.New("auth token not found")
var ErrMasterPasswordNotFound = errors.New("master password not found")

type SecretRepository struct {
	authTokenFileName      string
	masterPasswordFileName string
}

func NewSecretRepository(authTokenFilename string, masterPasswordFilename string) (*SecretRepository, error) {
	authTokenDir := filepath.Dir(authTokenFilename)
	masterPassDir := filepath.Dir(masterPasswordFilename)
	err := mkdir(authTokenDir)
	if err != nil {
		return nil, err
	}
	err = mkdir(masterPassDir)
	if err != nil {
		return nil, err
	}

	return &SecretRepository{authTokenFileName: authTokenFilename, masterPasswordFileName: masterPasswordFilename}, nil
}

func (s SecretRepository) SaveAuthToken(token string) error {
	file, err := os.OpenFile(s.authTokenFileName, os.O_WRONLY|os.O_CREATE, 0666)
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
	file, err := os.OpenFile(s.authTokenFileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", ErrTokenNotFound
	}

	return string(content), nil
}

func (s SecretRepository) SaveMasterPassword(pass string) error {
	file, err := os.OpenFile(s.masterPasswordFileName, os.O_WRONLY|os.O_CREATE, 0666)
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
	file, err := os.OpenFile(s.masterPasswordFileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", ErrMasterPasswordNotFound
	}

	return string(content), nil
}

func mkdir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}
