package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/iHamsin/practicum-shortener-service/internal/util"
)

type APIPostBatchInsertHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

type requestBatchInsertJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type responseBatchInsertJSON struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *APIPostBatchInsertHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

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

	codePrefix := "/"
	baseURL, _ := url.ParseRequestURI(h.Cfg.HTTP.BaseURL)
	if len(baseURL.Path) > 0 {
		codePrefix = ""
	}

	// Unmarshal request json
	var links []requestBatchInsertJSON
	jsonError := json.Unmarshal(body, &links)
	if jsonError != nil {
		http.Error(res, jsonError.Error(), http.StatusBadRequest)
		return
	}

	// Insert into repo
	originalLinks := make([]string, len(links))
	for i, link := range links {
		originalLinks[i] = link.OriginalURL
	}
	ctx := req.Context()
	UUID, _ := ctx.Value(middlewares.RequestUUIDKey{}).(string)
	shortLinks, repoError := h.Repo.BatchInsertLink(req.Context(), originalLinks, UUID)
	if repoError != nil {
		http.Error(res, repoError.Error(), http.StatusBadRequest)
		return
	}

	// Fill responce json
	results := make([]responseBatchInsertJSON, len(links))
	for i, link := range links {
		results[i] = responseBatchInsertJSON{
			CorrelationID: link.CorrelationID,
			ShortURL:      fmt.Sprintf("%s%s%s", h.Cfg.HTTP.BaseURL, codePrefix, shortLinks[i]),
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	error := json.NewEncoder(res).Encode(results)
	if error != nil {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	}
}
