package handlers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *PostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	if error != nil && !errors.Is(error, repositories.ErrDublicateOriginalLink) {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	} else {
		if errors.Is(error, repositories.ErrDublicateOriginalLink) {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}
		_, error := res.Write([]byte(fmt.Sprintf("%s%s%s", h.Cfg.HTTP.BaseURL, codePrefix, code)))
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}
	}
}
