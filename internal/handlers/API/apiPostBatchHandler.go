package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/sirupsen/logrus"
)

type APIPostBatchHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

type requestBatchJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type responseBatchJSON struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *APIPostBatchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			logrus.Debug("Error with gzip")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
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
	var links []requestBatchJSON
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
	UUID, _ := ctx.Value("UUID").(string)
	shortLinks, repoError := h.Repo.BatchInsertLink(req.Context(), originalLinks, UUID)
	if repoError != nil {
		http.Error(res, repoError.Error(), http.StatusBadRequest)
		return
	}

	// Fill responce json
	results := make([]responseBatchJSON, len(links))
	for i, link := range links {
		results[i] = responseBatchJSON{
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
