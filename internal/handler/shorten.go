package handler

import (
	"io"
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get(`Content-Type`) != `text/plain` {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalLink, err := io.ReadAll(r.Body)
	if err != nil || string(originalLink) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink := service.CreateShortLink(string(originalLink))

	w.Header().Add(`Content-Type`, `text/plain`)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortedLink))
}
