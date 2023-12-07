package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func Init(repository repositories.Repository, cfg config.Config) chi.Router {
	var router = chi.NewRouter()
	router.Use(middlewares.WithLoggingMiddleWare)
	router.Use(middlewares.WithCompressionResponse)

	RegisterAPIRoutes(router, repository, cfg)
	RegisterPublicRoutes(router, repository, cfg)
	return router
}
