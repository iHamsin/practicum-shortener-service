package main

import (
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"

	"github.com/go-chi/chi/v5"
)

func main() {
	// хранилище пока в памяти
	repository := repositories.NewLinksRepoRAM(make(map[string]string))

	router := chi.NewRouter()

	router.Post("/", handlers.InsertHandler(repository))
	router.Get("/{linkCode}", handlers.GetHandler(repository))

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
