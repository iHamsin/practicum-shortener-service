package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/iHamsin/practicum-shortener-service/internal/util"
	"github.com/sirupsen/logrus"
)

type APIUserDeleteURLSHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *APIUserDeleteURLSHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	UUID, _ := ctx.Value(middlewares.RequestUUIDKey{}).(string)

	var reader io.Reader

	reader, zipError := util.UnzipRequestBody(req)
	if zipError != nil {
		http.Error(res, zipError.Error(), http.StatusBadRequest)
		return
	}

	body, ioError := io.ReadAll(reader)
	if ioError != nil {
		http.Error(res, ioError.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal request json
	var links []string
	jsonError := json.Unmarshal(body, &links)
	if jsonError != nil {
		http.Error(res, jsonError.Error(), http.StatusBadRequest)
		return
	}

	_, batchDeleteError := h.Repo.BatchDeleteLink(ctx, links, UUID)

	if batchDeleteError != nil {
		logrus.Error(batchDeleteError.Error())
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
