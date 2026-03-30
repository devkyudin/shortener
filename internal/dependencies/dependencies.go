package dependencies

import (
	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler/getlink"
	"github.com/devkyudin/shortener/internal/handler/shortenjson"
	"github.com/devkyudin/shortener/internal/handler/shortenplaintext"
	"github.com/devkyudin/shortener/internal/logger"
	"github.com/devkyudin/shortener/internal/middleware"
	"github.com/devkyudin/shortener/internal/model"
	"github.com/devkyudin/shortener/internal/repository"
	shortener_router2 "github.com/devkyudin/shortener/internal/router"
	"github.com/devkyudin/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	LinksRepository         *repository.CodedLinksRepository
	Config                  *config.Config
	URLService              *service.URLService
	GetLinkHandler          *getlink.GetLinkHandler
	ShortenPlainTextHandler *shortenplaintext.ShortenPlainTextHandler
	ShortenJSONHandler      *shortenjson.ShortenJSONHandler
	Router                  *chi.Router
	LogContainer            *logger.Container
}

func GetDependencies() *Dependencies {
	cfg := config.GetConfig()
	logContainer := logger.NewLoggerContainer()
	alphabet := model.NewAlphabet([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	linksRepository, err := repository.NewCodedLinksRepository(cfg)
	if err != nil {
		logContainer.Logger.Error("ошибка при инициализации репозитория", "err", err)
		panic(err)
	}
	urlService := service.NewURLService(alphabet, linksRepository, cfg)
	getLinkHandler := getlink.NewGetLinkHandler(urlService)
	shortenPlainTextHandler := shortenplaintext.NewShortenPlainTextHandler(urlService)
	shortenJSONHandler := shortenjson.NewShortenJSONHandler(urlService)
	loggingMiddleware := middleware.NewLoggingMiddleware(logContainer)
	compressionMiddleware := middleware.NewCompressionMiddleware(logContainer, map[string]struct{}{"text/plain": {}, "application/json": {}})
	router := shortener_router2.GetRouter(shortenPlainTextHandler, shortenJSONHandler, getLinkHandler, loggingMiddleware, compressionMiddleware)
	return &Dependencies{
		LinksRepository:         linksRepository,
		Config:                  cfg,
		URLService:              urlService,
		GetLinkHandler:          getLinkHandler,
		ShortenPlainTextHandler: shortenPlainTextHandler,
		ShortenJSONHandler:      shortenJSONHandler,
		Router:                  router,
		LogContainer:            logContainer,
	}
}
