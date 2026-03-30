package ping

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/repository"
)

type PingHandler struct {
	r *repository.CodedLinksDbRepository
}

func NewPingHandler(repository *repository.CodedLinksDbRepository) *PingHandler {
	return &PingHandler{r: repository}
}

func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.r.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
