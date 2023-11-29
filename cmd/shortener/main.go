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

func main() {

	lvl, gotLvl := os.LookupEnv("LOG_LEVEL")
	if !gotLvl {
		lvl = "debug"
	}

	dbFile, gotDBFile := os.LookupEnv("FILE_STORAGE_PATH")
	if !gotDBFile {
		dbFile = "/tmp/short-url-db.json"
	}

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)

	cfg := &config.Config{}

	flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "HTTP server addr. Default: localhost:8080")
	flag.StringVar(&cfg.HTTP.BaseURL, "b", "http://localhost:8080", "Short link BaseURL. Default: http://localhost:8080")
	flag.StringVar(&cfg.HTTP.DBFile, "f", dbFile, "DB file. Default: /tmp/short-url-db.json")
	flag.Parse()

	configError := env.Parse(&cfg.HTTP)
	if configError != nil {
		logrus.Error(configError)
	}

	var repository repositories.Repository
	if cfg.HTTP.DBFile == "" {
		logrus.Debug("DB in RAM")
		repository = repositories.NewLinksRepoRAM(make(map[string]string))
	} else {
		logrus.Debug("DB in file", cfg.HTTP.DBFile)

		file, fileOpenError := os.OpenFile(cfg.HTTP.DBFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if fileOpenError != nil {
			logrus.Error(fileOpenError)
		}
		defer file.Close()

		var fileRepoError error
		repository, fileRepoError = repositories.NewLinksRepoFile(*file)
		if fileRepoError != nil {
			logrus.Error(fileRepoError)
		}
	}
	// хранилище пока в памяти

	router := chi.NewRouter()

	postHandler := &handlers.PostHandler{Repo: repository, Cfg: *cfg}
	getHandler := &handlers.GetHandler{Repo: repository, Cfg: *cfg}
	apiPostHandler := &handlers.APIPostHandler{Repo: repository, Cfg: *cfg}

	// Logger from Chi, too easy, will write custom with Logrus
	// router.Use(middleware.Logger)
	// router.Use(middleware.Recoverer)

	router.Use(middlewares.WithLoggingMiddleWare)
	router.Use(middlewares.WithCompressionResponse)

	router.Post("/", postHandler.ServeHTTP)
	router.Get("/{linkCode}", getHandler.ServeHTTP)

	router.Post("/api/shorten", apiPostHandler.ServeHTTP)

	logrus.Debug("WebServer started")

	fmt.Println("WebServer started at " + cfg.HTTP.Addr)
	fmt.Println("Short link BaseURL: " + cfg.HTTP.BaseURL)

	serverError := http.ListenAndServe(cfg.HTTP.Addr, router)
	if serverError != nil {
		logrus.Error(serverError)
	}

}
