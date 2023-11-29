package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func TestStatusHandler(t *testing.T) {
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
			name: "positive test #1",
			want: want{
				postCode:          201,
				getCode:           304,
				postBody:          "http://ok.kz",
				bodySize:          30,
				checkResponceBody: false,
				responceBody:      "",
				httpAddr:          "localhost:8080",
				httpBaseURL:       "http://localhost:8080",
			},
		},
		{
			name: "positive test #2",
			want: want{
				postCode:          400,
				getCode:           400,
				postBody:          "blablabla",
				bodySize:          43,
				checkResponceBody: true,
				responceBody:      "parse \"blablabla\": invalid URI for request\n",
				httpAddr:          "localhost:9090",
				httpBaseURL:       "http://localhost:9090/prefix-",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repository := repositories.NewLinksRepoRAM(make(map[string]string))
			cfg := new(config.Config)
			cfg.HTTP.Addr = test.want.httpAddr
			cfg.HTTP.BaseURL = test.want.httpBaseURL

			postHandler := PostHandler{Repo: repository, Cfg: *cfg}
			getHandler := GetHandler{Repo: repository, Cfg: *cfg}

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.want.postBody))
			w := httptest.NewRecorder()
			postHandler.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			resBody, _ := io.ReadAll(res.Body)
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

			request = httptest.NewRequest(http.MethodGet, string(resBody), nil)
			w = httptest.NewRecorder()
			getHandler.ServeHTTP(w, request)
			res = w.Result()
			defer res.Body.Close()
			// проверяем возврат линка по сохраненному коду
			assert.Equal(t, res.Header.Get("Location"), test.want.postBody)
		})
	}
}
