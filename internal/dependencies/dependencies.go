package dependencies

import (
	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler/getlink"
	"github.com/devkyudin/shortener/internal/handler/shortenjson"
	"github.com/devkyudin/shortener/internal/handler/shortenplaintext"
	"github.com/devkyudin/shortener/internal/logger"
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/devkyudin/shortener/internal/repository"
	shortener_router2 "github.com/devkyudin/shortener/internal/router"
	"github.com/devkyudin/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	LinksRepository         repository.LinksRepository
	Config                  config.Config
	URLService              service.URLService
	GetLinkHandler          getlink.GetLinkHandler
	ShortenPlainTextHandler shortenplaintext.ShortenPlainTextHandler
	ShortenJSONHandler      shortenjson.ShortenJSONHandler
	Router                  chi.Router
	LogContainer            logger.Container
}

func GetDependencies() *Dependencies {
	cfg := config.GetConfig()
	lr := repository.NewLinksRepository()
	s := service.NewURLService(lr, cfg)
	glh := getlink.NewGetLinkHandler(s)
	shp := shortenplaintext.NewShortenPlainTextHandler(s)
	shj := shortenjson.NewShortenJSONHandler(s)
	logContainer := logger.NewLoggerContainer()
	lm := middleware.NewLoggingMiddleware(logContainer)
	ch := middleware.NewCompressionMiddleware(logContainer, map[string]struct{}{"text/plain": {}, "application/json": {}})
	router := shortener_router2.GetRouter(shp, shj, glh, lm, ch)
	return &Dependencies{
		LinksRepository:         *lr,
		Config:                  *cfg,
		URLService:              *s,
		GetLinkHandler:          *glh,
		ShortenPlainTextHandler: *shp,
		ShortenJSONHandler:      *shj,
		Router:                  router,
		LogContainer:            *logContainer,
	}
}
