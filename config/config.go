package config

type (
	// Config -.
	Config struct {
		HTTP
	}

	// HTTP -.
	HTTP struct {
		Addr    string `env:"SERVER_ADDRESS"`
		BaseURL string `env:"BASE_URL"`
	}
)
