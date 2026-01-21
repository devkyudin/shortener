package handler

import (
	"io"
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

func Shorten(w http.ResponseWriter, r *http.Request) {

	originalLink, err := io.ReadAll(r.Body)
	if err != nil || string(originalLink) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink := service.CreateShortLink(string(originalLink))

	w.Header().Add(`Content-Type`, `text/plain`)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedLink))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
