package handler

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

func GetLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(`id`)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := service.GetFullLink(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add(`Content-Type`, `text/plain`)
	w.Header().Add(`Location`, result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
