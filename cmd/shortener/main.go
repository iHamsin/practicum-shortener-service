package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := new(config.Config)
	flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "HTTP server addr. Default: localhost:8080")
	flag.StringVar(&cfg.HTTP.BaseURL, "b", "http://localhost:8080", "Short link BaseURL. Default: http://localhost:8080")
	flag.Parse()

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
