package e2e

type Config struct {
	ServerBin           string `env:"SERVER_BIN"`
	ServerRunAddress    string `env:"SERVER_RUN_ADDRESS"`
	ServerDatabaseURI   string `env:"SERVER_DATABASE_URI"`
	ServerPublicAddress string `env:"SERVER_PUBLIC_ADDRESS"`

	ClientBin         string `env:"CLIENT_BIN"`
	ClientDataDirName string `env:"CLIENT_DATA_DIR_NAME"`
}

func NewConfig() *Config {
	return &Config{}
}
