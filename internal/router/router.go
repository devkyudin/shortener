package router

import (
	"github.com/devkyudin/shortener/internal/handler/getlink"
	"github.com/devkyudin/shortener/internal/handler/ping"
	"github.com/devkyudin/shortener/internal/handler/shortenjson"
	"github.com/devkyudin/shortener/internal/handler/shortenplaintext"
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRouter(shp *shortenplaintext.ShortenPlainTextHandler, shj *shortenjson.ShortenJSONHandler, gl *getlink.GetLinkHandler, p *ping.PingHandler, lm *middleware.LoggingMiddleware, cm *middleware.CompressionMiddleware) *chi.Router {
	r := chi.NewRouter()
	generalRouter := r.With(lm.WithLogging).With(cm.WithCompression)
	generalRouter.With(
		middleware.RequireContentType("text/plain"),
	).Post("/", shp.Shorten)

	generalRouter.With(
		middleware.RequireContentType("application/json"),
	).Post("/api/shorten", shj.Shorten)

	generalRouter.Get(`/{id}`, gl.GetLink)
	generalRouter.Get("/ping", p.Ping)
	result := chi.Router(r)
	return &result
}
