package config

const (
	defaultLoginFile = "./.data/entries/logins"
	defaultCardFile  = "./.data/entries/cards"
	defaultTextFile  = "./.data/entries/texts"
	defaultBinFile   = "./.data/entries/binaries"

	defaultAuthTokenFilename      = "./.data/secret/.token"
	defaultMasterPasswordFilename = "./.data/secret/.pass"
)

// Config missing godoc.
type Config struct {
	ServerAddress          string `env:"SERVER_ADDRESS"`
	LogLevel               string `env:"LOG_LEVEL"`
	LoginFilename          string `env:"LOGIN_FILENAME"`
	CardFilename           string `env:"CARD_FILENAME"`
	TextFilename           string `env:"TEXT_FILENAME"`
	BinFilename            string `env:"BIN_FILENAME"`
	AuthTokenFilename      string `env:"AUTH_TOKEN_FILENAME"`
	MasterPasswordFilename string `env:"MASTER_PASSWORD_FILENAME"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel:               "info",
		LoginFilename:          defaultLoginFile,
		CardFilename:           defaultCardFile,
		TextFilename:           defaultTextFile,
		BinFilename:            defaultBinFile,
		AuthTokenFilename:      defaultAuthTokenFilename,
		MasterPasswordFilename: defaultMasterPasswordFilename,
	}
}
