package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

type GetHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *GetHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// // разбиваем Path на сегменты
	// URLsegments := strings.Split(req.URL.Path, "/")
	// // если сегментов более одного, возвращаем ошибку
	// if len(URLsegments) != 2 {
	// 	http.Error(res, "Bad request", http.StatusBadRequest)
	// 	return
	// }
	// linkCode := URLsegments[1]

	// получение кода через Chi
	// отваливаются тесты если использовать chi.URLParam
	// linkCode := chi.URLParam(req, "linkCode")

	codePrefix := "/"
	baseURL, _ := url.ParseRequestURI(h.Cfg.HTTP.BaseURL)
	if len(baseURL.Path) > 0 {
		codePrefix = ""
	}

	linkCode := strings.TrimPrefix("http://"+string(req.Host+req.URL.Path), h.Cfg.HTTP.BaseURL+codePrefix)
	// logrus.Debug("New GET request with short code: " + linkCode)

	link, error := h.Repo.GetLinkByCode(req.Context(), linkCode)
	if error != nil {
		if error.Error() == "link gone" {
			http.Error(res, error.Error(), http.StatusGone)
		} else {
			http.Error(res, error.Error(), http.StatusBadRequest)
		}
		// logrus.Debug(error.Error())
	} else {
		res.Header().Set("Location", link)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
