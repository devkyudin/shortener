package service

import (
	"errors"
	"sync"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/model"
	"github.com/devkyudin/shortener/internal/repository"
)

type URLService struct {
	linksRepository *repository.CodedLinksFileRepository
	alphabet        *model.Alphabet
	cfg             *config.Config
	mutex           *sync.Mutex
}

func NewURLService(alphabet *model.Alphabet, linksRepository *repository.CodedLinksFileRepository, cfg *config.Config) *URLService {
	return &URLService{
		linksRepository: linksRepository,
		alphabet:        alphabet,
		cfg:             cfg,
		mutex:           &sync.Mutex{},
	}
}

func (s *URLService) CreateShortLink(originalURL string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	codedLink, isOk := s.linksRepository.GetByOriginalURL(originalURL)
	if isOk {
		return s.cfg.ShortLinkAddress.String() + codedLink.ShortURL, nil
	}

	id := s.linksRepository.GetUniqueID()
	codedLink = model.NewCodedLink(id, originalURL, s.alphabet)
	err := s.linksRepository.CreateCodedLink(codedLink)
	if err != nil {
		return "", errors.New("не удалось сохранить ссылку: " + err.Error())
	}
	shortedLink := s.cfg.ShortLinkAddress.String() + codedLink.ShortURL
	return shortedLink, nil
}

func (s *URLService) GetFullLink(shortURL string) (string, error) {
	codedLink, ok := s.linksRepository.GetByShortURL(shortURL)
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return codedLink.OriginalURL, nil
}
