package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func RegisterPublicRoutes(router chi.Router, repository repositories.Repository, cfg config.Config) {
	postHandler := &handlers.PostHandler{Repo: repository, Cfg: &cfg}
	getHandler := &handlers.GetHandler{Repo: repository, Cfg: &cfg}
	getDbPingHandler := &handlers.GetDBPingHandler{Repo: repository, Cfg: &cfg}

	router.Post("/", postHandler.ServeHTTP)
	router.Get("/{linkCode}", getHandler.ServeHTTP)
	router.Get("/ping", getDbPingHandler.ServeHTTP)
}
