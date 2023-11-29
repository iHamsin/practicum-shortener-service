package config

type Config struct {
	HTTP
}

// HTTP -.
type HTTP struct {
	Addr    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
	DBFile  string `env:"FILE_STORAGE_PATH"`
}
