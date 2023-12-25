package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iHamsin/practicum-shortener-service/config"
	handlers "github.com/iHamsin/practicum-shortener-service/internal/handlers/API"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func RegisterAPIRoutes(router chi.Router, repository repositories.Repository, cfg config.Config) {
	apiPostHandler := &handlers.APIPostHandler{Repo: repository, Cfg: &cfg}
	APIPostBatchInsertHandler := &handlers.APIPostBatchInsertHandler{Repo: repository, Cfg: &cfg}
	APIUserGetURLSHandler := &handlers.APIUserGetURLSHandler{Repo: repository, Cfg: &cfg}
	APIUserDeleteURLSHandler := &handlers.APIUserDeleteURLSHandler{Repo: repository, Cfg: &cfg}

	router.Post("/api/shorten", apiPostHandler.ServeHTTP)
	router.Post("/api/shorten/batch", APIPostBatchInsertHandler.ServeHTTP)
	router.Get("/api/user/urls", APIUserGetURLSHandler.ServeHTTP)
	router.Delete("/api/user/urls", APIUserDeleteURLSHandler.ServeHTTP)
}
