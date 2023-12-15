package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func TestStatusHandlerGzipJsonBatch(t *testing.T) {

	t.Run("positive batch JSON test", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080/addon/"
		cfg.Repository.ShortCodeLength = 8

		var repository, _ = repositories.Init(cfg)

		postHandler := APIPostBatchHandler{Repo: repository, Cfg: cfg}

		mcPostBody := []requestBatchJSON{{
			CorrelationID: "1",
			OriginalURL:   "https://practicum1.yandex.ru",
		}, {
			CorrelationID: "2",
			OriginalURL:   "https://practicum2.yandex.ru",
		}, {
			CorrelationID: "3",
			OriginalURL:   "https://practicum2.yandex.ru",
		}}

		body, _ := json.Marshal(mcPostBody)

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(body))
		if err != nil {
			fmt.Println(err)
		}

		zb.Close()

		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", buf)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")
		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(string(resBody))

		// проверяем код ответа
		assert.Equal(t, 201, res.StatusCode)
	})

	t.Run("negtive batch JSON test - broken Gzip", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080"
		cfg.Repository.ShortCodeLength = 8

		var repository, _ = repositories.Init(cfg)

		postHandler := APIPostBatchHandler{Repo: repository, Cfg: cfg}

		mcPostBody := []requestBatchJSON{{
			CorrelationID: "1",
			OriginalURL:   "brokenLink",
		}}
		body, _ := json.Marshal(mcPostBody)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")
		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(string(resBody))

		// проверяем код ответа
		assert.Equal(t, 500, res.StatusCode)
	})

	t.Run("negtive batch JSON test - broken Json", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080"
		cfg.Repository.ShortCodeLength = 8

		var repository, _ = repositories.Init(cfg)

		postHandler := APIPostBatchHandler{Repo: repository, Cfg: cfg}

		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader("brokenJson"))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(string(resBody))

		// проверяем код ответа
		assert.Equal(t, 400, res.StatusCode)
	})

}
