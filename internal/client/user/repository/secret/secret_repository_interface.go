package secret

//go:generate mockgen -source=secret_repository_interface.go -destination=mock_secret_repository/mock_secret_repository.go -package=mock_secret_repository
type SecretRepositoryInterface interface {
	SaveAuthToken(token string) error
	GetAuthToken() (string, error)

	SaveMasterPassword(pass string) error
	GetMasterPassword() (string, error)
}
