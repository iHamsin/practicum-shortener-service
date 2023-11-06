package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func InsertHandler(repo repositories.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)

		// парсим URL @todo надо найти лучше способ валидации URL
		_, error := url.ParseRequestURI(string(body))
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}

		// сохраняем линк
		code, error := repo.Insert(string(body))
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080/" + code))
		return
	}
}

func GetHandler(repo repositories.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// // разбиваем Path на сегменты
		// URLsegments := strings.Split(req.URL.Path, "/")
		// // если сегментов более одного, возвращаем ошибку
		// if len(URLsegments) != 2 {
		// 	http.Error(res, "Bad request", http.StatusBadRequest)
		// 	return
		// }

		linkCode := chi.URLParam(req, "linkCode")

		link, error := repo.GetByCode(linkCode)
		if error != nil {
			http.Error(res, error.Error(), http.StatusBadRequest)
			return
		} else {
			res.Header().Set("Location", link)
			res.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}
}
