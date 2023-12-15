package handlers

import (
	"bytes"
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

func TestStatusHandlerGzipJsonBatch(t *testing.T) {

	t.Run("test batch", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080"
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

		// mcPostBody := map[string]interface{}{
		// 	"url": "https://practicum.yandex.ru",
		// }

		body, _ := json.Marshal(mcPostBody)

		fmt.Println("----")
		fmt.Println(string(body))

		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(string(resBody))

		// проверяем код ответа
		assert.Equal(t, 201, res.StatusCode)
	})

}
