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

type APIPostHandler struct {
	Repo repositories.Repository
	Cfg  config.Config
}

type requestJSON struct {
	URL string
}

type responseJSON struct {
	Result string `json:"result"`
}

func (h *APIPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

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

	var reqJSON requestJSON
	jsonError := json.Unmarshal(body, &reqJSON)
	if jsonError != nil {
		http.Error(res, jsonError.Error(), http.StatusBadRequest)
		return
	}

	// парсим URL @todo надо найти лучше способ валидации URL
	_, error := url.ParseRequestURI(reqJSON.URL)

	if error != nil {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	}

	codePrefix := "/"
	baseURL, _ := url.ParseRequestURI(h.Cfg.HTTP.BaseURL)
	if len(baseURL.Path) > 0 {
		codePrefix = ""
	}

	// сохраняем линк
	code, error := h.Repo.Insert(reqJSON.URL)

	if error != nil {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	} else {
		var resJson responseJSON
		resJson.Result = fmt.Sprintf("%s%s%s", h.Cfg.HTTP.BaseURL, codePrefix, code)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		error := json.NewEncoder(res).Encode(resJson)
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}
	}
}
