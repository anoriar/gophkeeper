package config

const (
	defaultDataDirName = "./.data"
	defaultLoginFile   = "/entries/logins"
	defaultCardFile    = "/entries/cards"
	defaultTextFile    = "/entries/texts"
	defaultBinFile     = "/entries/binaries"

	defaultAuthTokenFilename      = "/secret/.token"
	defaultMasterPasswordFilename = "/secret/.pass"
)

// Config missing godoc.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
	DataDirName   string `env:"DATA_DIRNAME"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel:    "info",
		DataDirName: defaultDataDirName,
	}
}

func (cnf *Config) GetAuthTokenFilename() string {
	return cnf.DataDirName + defaultAuthTokenFilename
}

func (cnf *Config) GetMasterPasswordFilename() string {
	return cnf.DataDirName + defaultMasterPasswordFilename
}

func (cnf *Config) GetLoginFilename() string {
	return cnf.DataDirName + defaultLoginFile
}

func (cnf *Config) GetCardFilename() string {
	return cnf.DataDirName + defaultCardFile
}

func (cnf *Config) GetTextFilename() string {
	return cnf.DataDirName + defaultTextFile
}

func (cnf *Config) GetBinFilename() string {
	return cnf.DataDirName + defaultBinFile
}
