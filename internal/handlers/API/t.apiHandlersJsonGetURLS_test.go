package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/iHamsin/practicum-shortener-service/config"
	handlers "github.com/iHamsin/practicum-shortener-service/internal/handlers/public"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func TestGetUserURLS(t *testing.T) {

	t.Run("positive batch insert JSON test", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080/addon/"

		cfg.Repository.ShortCodeLength = 8
		cfg.Repository.DatabaseDSN = "postgres://yp:passw0rd@127.0.0.1:5432/postgres?sslmode=disable"

		var repository, repoError = repositories.Init(cfg)
		if repoError != nil {
			logrus.Error(repoError)
			return
		} else {
			defer repository.Close()
		}

		// создаем хэндлер
		postHandler := APIPostBatchInsertHandler{Repo: repository, Cfg: cfg}

		mcPostBody := []requestBatchInsertJSON{{
			CorrelationID: "1",
			OriginalURL:   "https://practicum1.yandex.ru",
		}, {
			CorrelationID: "2",
			OriginalURL:   "https://practicum2.yandex.ru",
		}, {
			CorrelationID: "3",
			OriginalURL:   "https://practicum3.yandex.ru",
		}}

		body, _ := json.Marshal(mcPostBody)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(body)))
		request.Header.Set("Content-Type", "application/json")

		// добавляем куку
		var key = "1234567890123456"
		UUID := "123"
		request = request.WithContext(context.WithValue(request.Context(), middlewares.RequestUUIDKey{}, UUID))
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(UUID))
		cryptedNewUUID := h.Sum(nil)
		request = request.WithContext(context.WithValue(request.Context(), middlewares.RequestisNewUUIDKey{}, cryptedNewUUID))

		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(string(resBody))

		// проверяем код ответа
		assert.Equal(t, 201, res.StatusCode)
	})

	t.Run("positive get and delete user URLS test", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080/addon/"

		cfg.Repository.ShortCodeLength = 8
		cfg.Repository.DatabaseDSN = "postgres://yp:passw0rd@127.0.0.1:5432/postgres?sslmode=disable"

		var repository, repoError = repositories.Init(cfg)
		if repoError != nil {
			logrus.Error(repoError)
			return
		} else {
			defer repository.Close()
		}

		// создаем хэндлер
		postHandler := APIUserGetURLSHandler{Repo: repository, Cfg: cfg}
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		// добавляем куку
		var key = "1234567890123456"
		UUID := "123"
		request = request.WithContext(context.WithValue(request.Context(), middlewares.RequestUUIDKey{}, UUID))
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(UUID))
		cryptedNewUUID := h.Sum(nil)
		request = request.WithContext(context.WithValue(request.Context(), middlewares.RequestisNewUUIDKey{}, cryptedNewUUID))

		w := httptest.NewRecorder()
		postHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		// проверяем код ответа
		assert.Equal(t, 200, res.StatusCode)

		//
		// массовое удаление
		//

		var links []responseBatchInsertJSON
		jsonError := json.Unmarshal(resBody, &links)
		assert.Equal(t, nil, jsonError)

		var linksToDelete []string
		for _, link := range links {
			linksToDelete = append(linksToDelete, strings.ReplaceAll(link.ShortURL, cfg.HTTP.BaseURL, ""))
		}
		fmt.Println(linksToDelete)

		body, _ := json.Marshal(linksToDelete)

		// создаем хэндлер
		deleteHandler := APIUserDeleteURLSHandler{Repo: repository, Cfg: cfg}

		deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
		deleteRequest.Header.Set("Content-Type", "application/json")
		deleteRequest.Header.Set("Cookie", "UUID="+UUID+"; UUIDSign="+string(cryptedNewUUID))

		w = httptest.NewRecorder()
		deleteHandler.ServeHTTP(w, deleteRequest)
		res = w.Result()
		defer res.Body.Close()
		resBody, _ = io.ReadAll(res.Body)

		fmt.Println(resBody)

		// проверяем код ответа
		assert.Equal(t, 202, res.StatusCode)

		// нужна пауза, иначе горутина не успевает удалить и следующий тест валится
		time.Sleep(3 * time.Second)

		//
		// получение удаленного линка
		//

		getHandler := handlers.GetHandler{Repo: repository, Cfg: cfg}
		checkDeleteRequest := httptest.NewRequest(http.MethodGet, string(links[0].ShortURL), nil)
		checkDeleteRequest.Header.Set("Cookie", "UUID="+UUID+"; UUIDSign="+string(cryptedNewUUID))
		w = httptest.NewRecorder()
		getHandler.ServeHTTP(w, checkDeleteRequest)
		res = w.Result()
		resBody, _ = io.ReadAll(res.Body)
		fmt.Println(resBody)
		defer res.Body.Close()
		// проверяем возврат линка по сохраненному коду
		assert.Equal(t, 410, res.StatusCode)
	})

	t.Run("negative delete user URLS test", func(t *testing.T) {

		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080/addon/"

		cfg.Repository.ShortCodeLength = 8
		cfg.Repository.DatabaseDSN = "postgres://yp:passw0rd@127.0.0.1:5432/postgres?sslmode=disable"

		var repository, repoError = repositories.Init(cfg)
		if repoError != nil {
			logrus.Error(repoError)
			return
		} else {
			defer repository.Close()
		}

		//
		// массовое удаление
		//

		// создаем хэндлер
		deleteHandler := APIUserDeleteURLSHandler{Repo: repository, Cfg: cfg}

		deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader([]byte("{}")))
		deleteRequest.Header.Set("Content-Type", "application/json")
		deleteRequest.Header.Set("Content-Encoding", "gzip")

		w := httptest.NewRecorder()
		deleteHandler.ServeHTTP(w, deleteRequest)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)

		fmt.Println(resBody)

		// проверяем код ответа
		assert.Equal(t, 500, res.StatusCode)
	})

}
