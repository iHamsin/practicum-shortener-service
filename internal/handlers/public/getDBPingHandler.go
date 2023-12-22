package handlers

import (
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/sirupsen/logrus"
)

type GetDBPingHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *GetDBPingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logrus.Debug("New DB ping request")
	if h.Repo == nil {
		logrus.Debug("DB check failed")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	checkError := h.Repo.Check()
	if checkError != nil {
		logrus.Debug("DB check failed")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
