package shortenjson

import (
	"encoding/json"
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

type ShortenJSONHandler struct {
	s *service.URLService
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func NewShortenJSONHandler(s *service.URLService) *ShortenJSONHandler {
	return &ShortenJSONHandler{s: s}
}

func (h *ShortenJSONHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	bodyRequest := ShortenRequest{}
	if err := decoder.Decode(&bodyRequest); err != nil || bodyRequest.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink := h.s.CreateShortLink(bodyRequest.URL)

	w.Header().Add(`Content-Type`, `application/json`)
	w.WriteHeader(http.StatusCreated)
	response := ShortenResponse{Result: shortedLink}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
