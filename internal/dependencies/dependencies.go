package dependencies

import (
	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler/getlink"
	"github.com/devkyudin/shortener/internal/handler/ping"
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
	LinksFileRepository     *repository.CodedLinksFileRepository
	LinksDbRepository       *repository.CodedLinksDbRepository
	Config                  *config.Config
	URLService              *service.URLService
	GetLinkHandler          *getlink.GetLinkHandler
	ShortenPlainTextHandler *shortenplaintext.ShortenPlainTextHandler
	ShortenJSONHandler      *shortenjson.ShortenJSONHandler
	PingHandler             *ping.PingHandler
	Router                  *chi.Router
	LogContainer            *logger.Container
}

func GetDependencies() *Dependencies {
	cfg := config.GetConfig()
	logContainer := logger.NewLoggerContainer()
	alphabet := model.NewAlphabet([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	linksFileRepository, err := repository.NewCodedLinksFileRepository(cfg)
	if err != nil {
		logContainer.Logger.Error("ошибка при инициализации файлового репозитория", "err", err)
		panic(err)
	}

	linksDbRepository, err := repository.NewCodedLinksDbRepository(cfg, logContainer)
	if err != nil {
		logContainer.Logger.Error("ошибка при инициализации БД", "err", err)
		panic(err)
	}
	urlService := service.NewURLService(alphabet, linksFileRepository, cfg)

	getLinkHandler := getlink.NewGetLinkHandler(urlService)
	shortenPlainTextHandler := shortenplaintext.NewShortenPlainTextHandler(urlService)
	shortenJSONHandler := shortenjson.NewShortenJSONHandler(urlService)
	pingHandler := ping.NewPingHandler(linksDbRepository)
	loggingMiddleware := middleware.NewLoggingMiddleware(logContainer)
	compressionMiddleware := middleware.NewCompressionMiddleware(logContainer, map[string]struct{}{"text/plain": {}, "application/json": {}})
	router := shortener_router2.GetRouter(shortenPlainTextHandler, shortenJSONHandler, getLinkHandler, pingHandler, loggingMiddleware, compressionMiddleware)
	return &Dependencies{
		LinksFileRepository:     linksFileRepository,
		LinksDbRepository:       linksDbRepository,
		Config:                  cfg,
		URLService:              urlService,
		GetLinkHandler:          getLinkHandler,
		ShortenPlainTextHandler: shortenPlainTextHandler,
		ShortenJSONHandler:      shortenJSONHandler,
		PingHandler:             pingHandler,
		Router:                  router,
		LogContainer:            logContainer,
	}
}
