package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func TestStatusHandlerPing(t *testing.T) {

	t.Run("positive pingDB test", func(t *testing.T) {
		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080"
		cfg.Repository.ShortCodeLength = 8
		var repository, _ = repositories.Init(cfg)
		pingHandler := GetDBPingHandler{Repo: repository, Cfg: cfg}
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		pingHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		fmt.Println(string(resBody))
		assert.Equal(t, 200, res.StatusCode)
	})
	t.Run("negative pingDB test - no DB", func(t *testing.T) {
		cfg := new(config.Config)
		cfg.HTTP.Addr = "localhost:8080"
		cfg.HTTP.BaseURL = "http://localhost:8080"
		cfg.Repository.DatabaseDSN = "postgres://user:password@localhost:5432/postgres?sslmode=disable"
		cfg.Repository.ShortCodeLength = 8
		var repository, _ = repositories.Init(cfg)
		pingHandler := GetDBPingHandler{Repo: repository, Cfg: cfg}
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		pingHandler.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		fmt.Println(string(resBody))
		assert.Equal(t, 500, res.StatusCode)
	})

}
