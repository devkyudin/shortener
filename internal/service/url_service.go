package service

import (
	"errors"
	"sync"
	"unicode/utf8"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/repository"
)

type URLService struct {
	lr           *repository.LinksRepository
	cfg          *config.Config
	nextID       int
	mutex        *sync.Mutex
	linkAlphabet []rune
	alphabetMap  map[rune]int
}

func NewURLService(lr *repository.LinksRepository, cfg *config.Config) *URLService {
	linkAlphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	alphabetMap := createMap(linkAlphabet)
	return &URLService{
		lr:           lr,
		cfg:          cfg,
		nextID:       10_000_000,
		mutex:        &sync.Mutex{},
		linkAlphabet: linkAlphabet,
		alphabetMap:  alphabetMap,
	}
}

func (s *URLService) CreateShortLink(originalLink string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, isOk := s.lr.GetByLink(originalLink)
	if isOk {
		// Если ссылка уже есть, возвращаем существующую короткую ссылку
		id, _ := s.lr.GetByLink(originalLink)
		return s.cfg.ShortLinkAddress.String() + s.toString(id)
	}

	id := s.getNewID()
	s.lr.CreateShortLink(originalLink, id)
	shortedLink := s.cfg.ShortLinkAddress.String() + s.toString(id)
	return shortedLink
}

func (s *URLService) getNewID() int {
	result := s.nextID
	s.nextID++
	return result
}

func (s *URLService) toString(id int) string {
	var result = ""
	alphabetLength := len(s.linkAlphabet)
	for id > 0 {
		code := id % alphabetLength
		result = result + string(s.linkAlphabet[code])
		id = id / alphabetLength
	}

	return reverse(result)
}

func (s *URLService) stringToID(src string) (int, error) {
	result := 0
	alphabetLength := len(s.linkAlphabet)
	reversed := reverse(src)
	runes := []rune(reversed)
	for i := 0; i < len(src); i++ {
		runeID, ok := s.alphabetMap[runes[i]]
		if !ok {
			return 0, errors.New(`битая ссылка`)
		}
		result += runeID * intPow(alphabetLength, i)
	}

	return result, nil
}

func createMap(alphabet []rune) map[rune]int {
	result := make(map[rune]int, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		result[alphabet[i]] = i
	}

	return result
}

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func (s *URLService) GetFullLink(codedID string) (string, error) {
	id, err := s.stringToID(codedID)
	if err != nil {
		return "", err
	}

	fullLink, ok := s.lr.GetByID(id)
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return fullLink, nil
}
