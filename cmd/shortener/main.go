package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := new(config.Config)

	configError := env.Parse(&cfg.HTTP)
	if configError != nil {
		log.Fatal(configError)
	}
	if cfg.HTTP.Addr == "" {
		flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "HTTP server addr. Default: localhost:8080")
		flag.Parse()
	}
	if cfg.HTTP.BaseURL == "" {
		flag.StringVar(&cfg.HTTP.BaseURL, "b", "http://localhost:8080", "Short link BaseURL. Default: http://localhost:8080")
		flag.Parse()
	}

	// хранилище пока в памяти
	repository := repositories.NewLinksRepoRAM(make(map[string]string))

	router := chi.NewRouter()

	router.Post("/", handlers.InsertHandler(repository, *cfg))
	router.Get("/{linkCode}", handlers.GetHandler(repository, *cfg))

	fmt.Println("WebServer started at " + cfg.HTTP.Addr)
	fmt.Println("Short link BaseURL: " + cfg.HTTP.BaseURL)

	err := http.ListenAndServe(cfg.HTTP.Addr, router)
	if err != nil {
		panic(err)
	}

}
