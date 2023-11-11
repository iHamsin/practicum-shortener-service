package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

type PostHandler struct {
	Repo repositories.Repository
	Cfg  config.Config
}

func (h *PostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	// парсим URL @todo надо найти лучше способ валидации URL
	_, error := url.ParseRequestURI(string(body))

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
	code, error := h.Repo.Insert(string(body))
	if error != nil {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	} else {
		res.WriteHeader(http.StatusCreated)
		_, error := res.Write([]byte(fmt.Sprintf("%s%s%s", h.Cfg.HTTP.BaseURL, codePrefix, code)))
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}
	}
}
