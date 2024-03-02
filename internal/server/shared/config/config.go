package config

// Config missing godoc.
type Config struct {
	RunAddress   string `env:"RUN_ADDRESS"`
	LogLevel     string `env:"LOG_LEVEL"`
	DatabaseURI  string `env:"DATABASE_URI"`
	JwtSecretKey string `env:"JWT_SECRET_KEY"`
}

// NewConfig missing godoc.
func NewConfig() *Config {
	return &Config{
		LogLevel:     "info",
		JwtSecretKey: "secret-key",
	}
}
