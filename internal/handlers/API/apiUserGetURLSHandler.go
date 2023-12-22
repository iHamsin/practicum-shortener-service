package handlers

import (
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

type APIUserGetURLSHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *APIUserGetURLSHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
