package shortenplaintext

import (
	"io"
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

type ShortenPlainTextHandler struct {
	s *service.URLService
}

func NewShortenPlainTextHandler(s *service.URLService) *ShortenPlainTextHandler {
	return &ShortenPlainTextHandler{s: s}
}

func (h *ShortenPlainTextHandler) Shorten(w http.ResponseWriter, r *http.Request) {

	originalLink, err := io.ReadAll(r.Body)
	if err != nil || string(originalLink) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink, err := h.s.CreateShortLink(string(originalLink))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(`Content-Type`, `text/plain`)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortedLink))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
