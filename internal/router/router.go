package router

import (
	"github.com/devkyudin/shortener/internal/handler/getlink"
	shortener_router2 "github.com/devkyudin/shortener/internal/handler/shorten"
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRouter(sh *shortener_router2.ShortenHandler, gl *getlink.GetLinkHandler) chi.Router {
	r := chi.NewRouter()
	r.With(
		middleware.RequireContentType("text/plain"),
	).Post("/", sh.Shorten)

	r.Get(`/{id}`, gl.GetLink)
	return r
}
