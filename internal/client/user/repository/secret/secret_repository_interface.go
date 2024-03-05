package secret

type SecretRepositoryInterface interface {
	SaveAuthToken(token string) error
	GetAuthToken() (string, error)

	SaveMasterPassword(pass string) error
	GetMasterPassword() (string, error)
}
