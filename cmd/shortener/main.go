package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/iHamsin/practicum-shortener-service/internal/routes"
	"github.com/sirupsen/logrus"
)

func main() {

	lvl, gotLvl := os.LookupEnv("LOG_LEVEL")
	if !gotLvl {
		lvl = "debug"
	}

	// logrus init
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)

	// config init
	cfg, configError := config.Init()
	if configError != nil {
		logrus.Error(configError)
	}

	var repository, repoError = repositories.Init(&cfg)
	if repoError != nil {
		logrus.Error(configError)
	}
	defer repository.Close()

	router := routes.Init(repository, cfg)
	serverError := http.ListenAndServe(cfg.HTTP.Addr, router)
	if serverError != nil {
		logrus.Error(serverError)
	}
	logrus.Debug("WebServer started")
	fmt.Println("WebServer started at " + cfg.HTTP.Addr)
	fmt.Println("Short link BaseURL: " + cfg.HTTP.BaseURL)
}
