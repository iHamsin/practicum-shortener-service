package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func RegisterAPIRoutes(router chi.Router, repository repositories.Repository, cfg config.Config) {
	apiPostHandler := &handlers.APIPostHandler{Repo: repository, Cfg: &cfg}

	router.Post("/api/shorten", apiPostHandler.ServeHTTP)
}
