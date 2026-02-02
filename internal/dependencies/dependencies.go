package dependencies

import (
	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler/getlink"
	"github.com/devkyudin/shortener/internal/handler/shorten"
	"github.com/devkyudin/shortener/internal/logger"
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/devkyudin/shortener/internal/repository"
	shortener_router2 "github.com/devkyudin/shortener/internal/router"
	"github.com/devkyudin/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	LinksRepository repository.LinksRepository
	Config          config.Config
	URLService      service.URLService
	GetLinkHandler  getlink.GetLinkHandler
	ShortenHandler  shorten.ShortenHandler
	Router          chi.Router
	LogContainer    logger.Container
}

func GetDependencies() *Dependencies {
	cfg := config.GetConfig()
	lr := repository.NewLinksRepository()
	s := service.NewURLService(lr, cfg)
	glh := getlink.NewGetLinkHandler(s)
	sh := shorten.NewShortenHandler(s)
	logContainer := logger.NewLoggerContainer()
	lm := middleware.NewLoggingMiddleware(logContainer)
	router := shortener_router2.GetRouter(sh, glh, lm)
	return &Dependencies{
		LinksRepository: *lr,
		Config:          *cfg,
		URLService:      *s,
		GetLinkHandler:  *glh,
		ShortenHandler:  *sh,
		Router:          router,
		LogContainer:    *logContainer,
	}
}
