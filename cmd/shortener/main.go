package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func main() {
	cfg := &config.Config{}

	flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "HTTP server addr. Default: localhost:8080")
	flag.StringVar(&cfg.HTTP.BaseURL, "b", "http://localhost:8080", "Short link BaseURL. Default: http://localhost:8080")
	flag.Parse()

	configError := env.Parse(&cfg.HTTP)
	if configError != nil {
		logrus.Error(configError)
	}

	// хранилище пока в памяти
	repository := repositories.NewLinksRepoRAM(make(map[string]string))

	router := chi.NewRouter()

	postHandler := &handlers.PostHandler{Repo: repository, Cfg: *cfg}
	getHandler := &handlers.GetHandler{Repo: repository, Cfg: *cfg}

	// Logger from Chi, too easy, will write custom with Logrus
	// router.Use(middleware.Logger)
	// router.Use(middleware.Recoverer)

	router.Use(middlewares.WithLoggingMiddleWare)

	router.Post("/", postHandler.ServeHTTP)
	router.Get("/{linkCode}", getHandler.ServeHTTP)

	logrus.Debug("WebServer started")

	fmt.Println("WebServer started at " + cfg.HTTP.Addr)
	fmt.Println("Short link BaseURL: " + cfg.HTTP.BaseURL)

	serverError := http.ListenAndServe(cfg.HTTP.Addr, router)
	if serverError != nil {
		logrus.Error(serverError)
	}

}
