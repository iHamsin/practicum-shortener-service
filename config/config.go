package config

type (
	// Config -.
	Config struct {
		HTTP
	}

	// HTTP -.
	HTTP struct {
		Addr    string
		BaseURL string
	}
)
