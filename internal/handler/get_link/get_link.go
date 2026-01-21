package get_link

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/service"
)

type GetLinkHandler struct {
	s *service.URLService
}

func NewGetLinkHandler(s *service.URLService) *GetLinkHandler {
	return &GetLinkHandler{s: s}
}

func (h *GetLinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(`id`)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.s.GetFullLink(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add(`Content-Type`, `text/plain`)
	w.Header().Add(`Location`, result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
