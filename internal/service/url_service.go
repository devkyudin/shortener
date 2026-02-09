package service

import (
	"errors"
	"sync"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/model"
	"github.com/devkyudin/shortener/internal/repository"
)

type URLService struct {
	linksRepository *repository.CodedLinksRepository
	alphabet        *model.Alphabet
	cfg             *config.Config
	mutex           *sync.Mutex
}

func NewURLService(alphabet *model.Alphabet, linksRepository *repository.CodedLinksRepository, cfg *config.Config) *URLService {
	return &URLService{
		linksRepository: linksRepository,
		alphabet:        alphabet,
		cfg:             cfg,
		mutex:           &sync.Mutex{},
	}
}

func (s *URLService) CreateShortLink(originalUrl string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	codedLink, isOk := s.linksRepository.GetByOriginalUrl(originalUrl)
	if isOk {
		return s.cfg.ShortLinkAddress.String() + codedLink.ShortUrl
	}

	id := s.linksRepository.GetUniqueID()
	codedLink = model.NewCodedLink(id, originalUrl, s.alphabet)
	s.linksRepository.CreateCodedLink(codedLink)
	shortedLink := s.cfg.ShortLinkAddress.String() + codedLink.ShortUrl
	return shortedLink
}

func (s *URLService) GetFullLink(shortUrl string) (string, error) {
	codedLink, ok := s.linksRepository.GetByShortUrl(shortUrl)
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return codedLink.OriginalUrl, nil
}
