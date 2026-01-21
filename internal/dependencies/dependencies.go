package dependencies

import (
	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler/get_link"
	"github.com/devkyudin/shortener/internal/handler/shorten"
	"github.com/devkyudin/shortener/internal/repository"
	shortener_router2 "github.com/devkyudin/shortener/internal/router"
	"github.com/devkyudin/shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	LinksRepository repository.LinksRepository
	Config          config.Config
	URLService      service.URLService
	GetLinkHandler  get_link.GetLinkHandler
	ShortenHandler  shorten.ShortenHandler
	Router          chi.Router
}

func GetDependencies() *Dependencies {
	cfg := config.GetConfig()
	lr := repository.NewLinksRepository()
	s := service.NewURLService(lr, cfg)
	glh := get_link.NewGetLinkHandler(s)
	sh := shorten.NewShortenHandler(s)
	router := shortener_router2.GetRouter(sh, glh)
	return &Dependencies{
		LinksRepository: *lr,
		Config:          *cfg,
		URLService:      *s,
		GetLinkHandler:  *glh,
		ShortenHandler:  *sh,
		Router:          router,
	}
}
