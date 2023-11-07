package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func InsertHandler(repo repositories.Repository, cfg config.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)

		// парсим URL @todo надо найти лучше способ валидации URL
		_, error := url.ParseRequestURI(string(body))

		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}

		codePrefix := "/"
		baseURL, _ := url.ParseRequestURI(cfg.HTTP.BaseURL)
		if len(baseURL.Path) > 0 {
			codePrefix = ""
		}

		// сохраняем линк
		code, error := repo.Insert(string(body))
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(cfg.HTTP.BaseURL + codePrefix + code))
		}
	}
}

func GetHandler(repo repositories.Repository, cfg config.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

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
		baseURL, _ := url.ParseRequestURI(cfg.HTTP.BaseURL)
		if len(baseURL.Path) > 0 {
			codePrefix = ""
		}

		fmt.Println("> GET request")
		linkCode := strings.TrimPrefix("http://"+string(req.Host+req.URL.Path), cfg.HTTP.BaseURL+codePrefix)
		fmt.Println("Code asked: " + linkCode)

		link, error := repo.GetByCode(linkCode)
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
		} else {
			res.Header().Set("Location", link)
			res.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
}
