package handler

import (
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func ShortenerRouter() chi.Router {
	r := chi.NewRouter()
	r.With(
		middleware.RequireContentType("text/plain"),
	).Post("/", Shorten)

	r.Get(`/{id}`, GetLink)
	return r
}
