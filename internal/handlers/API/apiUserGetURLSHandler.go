package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/iHamsin/practicum-shortener-service/internal/middlewares"
	"github.com/iHamsin/practicum-shortener-service/internal/repositories"
)

type APIUserGetURLSHandler struct {
	Repo repositories.Repository
	Cfg  *config.Config
}

func (h *APIUserGetURLSHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	UUID, _ := ctx.Value(middlewares.RequestUUIDKey{}).(string)
	isNewUUID, _ := ctx.Value(middlewares.RequestisNewUUIDKey{}).(bool)

	links, err := h.Repo.GetLinksByUUID(req.Context(), UUID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if isNewUUID {
		res.WriteHeader(http.StatusUnauthorized)
	} else if len(links) == 0 {
		res.WriteHeader(http.StatusNoContent)
	} else {
		res.WriteHeader(http.StatusOK)
	}

	codePrefix := "/"
	baseURL, _ := url.ParseRequestURI(h.Cfg.HTTP.BaseURL)
	if len(baseURL.Path) > 0 {
		codePrefix = ""
	}

	for i := range links {
		links[i].ShortURL = fmt.Sprintf("%s%s%s", h.Cfg.HTTP.BaseURL, codePrefix, links[i].ShortURL)
	}

	error := json.NewEncoder(res).Encode(links)
	if error != nil {
		http.Error(res, error.Error(), http.StatusBadRequest)
		return
	}
}
