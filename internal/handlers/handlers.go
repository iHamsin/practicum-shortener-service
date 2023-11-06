package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func MainHandler(repo repositories.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case "POST":
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
		case "GET":
			// разбиваем Path на сегменты
			URLsegments := strings.Split(req.URL.Path, "/")
			// если сегментов более одного, возвращаем ошибку
			if len(URLsegments) != 2 {
				http.Error(res, "Bad request", http.StatusBadRequest)
				return
			}

			link, error := repo.GetByCode(URLsegments[1])
			if error != nil {
				http.Error(res, error.Error(), http.StatusBadRequest)
				return
			} else {
				res.Header().Set("Location", link)
				res.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		default:
			http.Error(res, "Only GET and POST requests are allowed!", http.StatusBadRequest)
			return
		}
	}
}