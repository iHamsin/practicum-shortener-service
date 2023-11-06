package main

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

func randomString(n int) string {
	// словарь
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alphabetSize := len(alphabet)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}
	return sb.String()
}

func mainHandler(storage map[string]string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet { // обрабатываем GET запрос
			// разбиваем Path на сегменты
			URLsegments := strings.Split(req.URL.Path, "/")
			// если сегментов более одного, возвращаем ошибку
			if len(URLsegments) != 2 {
				http.Error(res, "Bad request", http.StatusBadRequest)
				return
			}
			// проверка наличия в хранилище
			_, URLfound := storage[URLsegments[1]]
			if URLfound {
				res.Header().Set("Location", storage[URLsegments[1]])
				res.WriteHeader(http.StatusTemporaryRedirect)
				return
			} else {
				http.Error(res, "Not Found", http.StatusBadRequest)
				return
			}
		} else if req.Method == http.MethodPost { // обрабатываем POST запрос - новый идентификатор сокращённого URL
			body, _ := io.ReadAll(req.Body)
			// парсим URL @todo надо найти лучше способ валидации URL
			_, err := url.ParseRequestURI(string(body))
			// при ошибке парсинг URL возвращаем 400 и описание ошики
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			// генерируем ключ и проверяем на наличие такого в хранилище
			var URLkey string
			var i = 0
			for {
				URLkey = randomString(8)
				_, keyUsed := storage[URLkey]
				if !keyUsed {
					break
				}
				i++
				// если за 20 попыток не сгенерировали униклаьный ключ - сдаемся
				if i > 20 {
					http.Error(res, "Cant generate uniq URL short code", http.StatusInternalServerError)
					return
				}
			}
			// записываем в хранилище
			storage[URLkey] = string(body)
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte("http://localhost:8080/" + URLkey))
			// @temp для проверки хранилища
			// for k, v := range storage {
			// 	res.Write([]byte(fmt.Sprintf("%s: %v\r\n", k, v)))
			// }
			return
		} else { // метод отличный от POST или GET
			http.Error(res, "Only GET and POST requests are allowed!", http.StatusBadRequest)
			return
		}
	}
}

func main() {
	// хранилище пока в памяти
	storage := make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainHandler(storage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
