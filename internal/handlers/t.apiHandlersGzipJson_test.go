package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func TestStatusHandlerJson(t *testing.T) {
	type want struct {
		postCode          int
		getCode           int
		postBody          string
		bodySize          int
		checkResponceBody bool
		responceBody      string
		httpAddr          string
		httpBaseURL       string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test json",
			want: want{
				postCode:          201,
				getCode:           304,
				postBody:          `{"url": "https://ok.kz"}`,
				bodySize:          44,
				checkResponceBody: false,
				responceBody:      "",
				httpAddr:          "localhost:8080/api/shorten",
				httpBaseURL:       "http://localhost:8080",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			cfg := new(config.Config)
			cfg.HTTP.Addr = test.want.httpAddr
			cfg.HTTP.BaseURL = test.want.httpBaseURL
			cfg.Repository.ShortCodeLength = 8

			var repository, _ = repositories.Init(cfg)

			postHandler := APIPostHandler{Repo: repository, Cfg: cfg}

			mcPostBody := map[string]interface{}{
				"url": "https://practicum.yandex.ru",
			}
			body, _ := json.Marshal(mcPostBody)

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(body))
			if err != nil {
				fmt.Println(err)
			}

			zb.Close()

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", buf)

			request.Header.Set("Content-Type", "application/json; charset=utf-8")
			request.Header.Set("Content-Encoding", "gzip")
			w := httptest.NewRecorder()
			postHandler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, _ := io.ReadAll(res.Body)

			fmt.Println("------")
			fmt.Println(buf)
			fmt.Println(string(resBody))

			// проверяем код ответа
			assert.Equal(t, test.want.postCode, res.StatusCode)
			// проверяем длину ответа, код рандомный, только так
			assert.Equal(t, len(resBody), test.want.bodySize)
			if test.want.checkResponceBody {
				// проверяем содержание ответа если там ошибка
				assert.Equal(t, string(resBody), test.want.responceBody)
			}

			// если была ошибка, выходим, проверять GET нет смысла
			if test.want.postCode == 400 {
				return
			}
		})
	}
}
