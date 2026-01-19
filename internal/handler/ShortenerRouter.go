package handler

import (
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func ShortenerRouter() chi.Router {
	r := chi.NewRouter()
	r.With(
		middleware.RequireMethod("POST"),
		middleware.RequireContentType("text/plain"),
		middleware.RequireNonEmptyBody(),
	).HandleFunc("/", Shorten)

	r.With(
		middleware.RequireMethod("GET"),
	).HandleFunc(`/{id}`, GetLink)
	return r
}
