package config

// Config missing godoc.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel: "info",
	}
}
