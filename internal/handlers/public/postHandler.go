package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
	"github.com/iHamsin/practicum-shortener-service/internal/util"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *PostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var reader io.Reader

	reader, zipError := util.UnzipRequestBody(req)
	if zipError != nil {
		http.Error(res, zipError.Error(), http.StatusInternalServerError)
		return
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
	ctx := req.Context()
	UUID, _ := ctx.Value(middlewares.RequestUUIDKey{}).(string)

	// сохраняем линк
	code, error := h.Repo.InsertLink(req.Context(), string(body), UUID)
	if error != nil && !errors.Is(error, repositories.ErrDublicateOriginalLink) {
		http.Error(res, error.Error(), http.StatusBadRequest)
		logrus.Error(error)
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
