package main

import (
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/internal/handlers"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

func main() {
	// хранилище пока в памяти
	repository := repositories.NewLinksRepoRAM(make(map[string]string))

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.MainHandler(repository))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
