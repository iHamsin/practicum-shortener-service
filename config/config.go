package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	HTTP
	Repository
}

// HTTP -.
type HTTP struct {
	Addr    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
	DBFile  string `env:"FILE_STORAGE_PATH"`
}

// Repository -.
type Repository struct {
	ShortCodeLength int    `env:"ShortCodeLength"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func Init() (Config, error) {
	cfg := &Config{}

	dbFile, gotDBFile := os.LookupEnv("FILE_STORAGE_PATH")
	if !gotDBFile {
		dbFile = "/tmp/short-url-db.json"
	}

	flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "HTTP server addr. Default: localhost:8080")
	flag.StringVar(&cfg.HTTP.BaseURL, "b", "http://localhost:8080", "Short link BaseURL. Default: http://localhost:8080")
	flag.StringVar(&cfg.HTTP.DBFile, "f", dbFile, "DB file. Example: /tmp/short-url-db.json")
	flag.IntVar(&cfg.Repository.ShortCodeLength, "l", 8, "Short code length. Default: 8")
	flag.StringVar(&cfg.Repository.DatabaseDSN, "d", "", "Database DSN. Example: host=localhost user=yp password=pa$$e0rd dbname=shortener sslmode=disable")
	flag.Parse()

	configError := env.Parse(&cfg.HTTP)

	if configError != nil {
		return *cfg, configError
	}

	configError = env.Parse(&cfg.Repository)

	if configError != nil {
		return *cfg, configError
	}

	return *cfg, nil
}
