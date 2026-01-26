package shorten

import (
	"io"
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

type ShortenHandler struct {
	s *service.URLService
}

func NewShortenHandler(s *service.URLService) *ShortenHandler {
	return &ShortenHandler{s: s}
}

func (h *ShortenHandler) Shorten(w http.ResponseWriter, r *http.Request) {

	originalLink, err := io.ReadAll(r.Body)
	if err != nil || string(originalLink) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink := h.s.CreateShortLink(string(originalLink))

	w.Header().Add(`Content-Type`, `text/plain`)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedLink))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
