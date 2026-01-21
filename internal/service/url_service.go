package service

import (
	"errors"
	"sync"
	"unicode/utf8"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/repository"
)

type URLService struct {
	lr  *repository.LinksRepository
	cfg *config.Config
}

func NewURLService(lr *repository.LinksRepository, cfg *config.Config) *URLService {
	return &URLService{lr, cfg}
}

var newShortLinkID = initialShortLinkID
var initialShortLinkID = 10_000_000
var mutex sync.Mutex
var shortLinkAlphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var alphabetMap map[rune]int

func (s *URLService) CreateShortLink(originalLink string) string {
	mutex.Lock()
	defer mutex.Unlock()
	_, isOk := s.lr.GetByLink(originalLink)
	if isOk {
		// Если ссылка уже есть, возвращаем существующую короткую ссылку
		id, _ := s.lr.GetByLink(originalLink)
		return s.cfg.ShortLinkAddress.String() + toString(id)
	}

	id := getNewID()
	s.lr.CreateShortLink(originalLink, id)
	shortedLink := s.cfg.ShortLinkAddress.String() + toString(id)
	return shortedLink
}

func getNewID() int {
	result := newShortLinkID
	newShortLinkID++
	return result
}

func toString(id int) string {
	var result = ""
	alphabetLength := len(shortLinkAlphabet)
	for id > 0 {
		code := id % alphabetLength
		result = result + string(shortLinkAlphabet[code])
		id = id / alphabetLength
	}

	return reverse(result)
}

func stringToID(src string) (int, error) {
	if alphabetMap == nil {
		alphabetMap = createMap()
	}
	result := 0
	alphabetLength := len(shortLinkAlphabet)
	reversed := reverse(src)
	runes := []rune(reversed)
	for i := 0; i < len(src); i++ {
		runeID, ok := alphabetMap[runes[i]]
		if !ok {
			return 0, errors.New(`битая ссылка`)
		}
		result += runeID * intPow(alphabetLength, i)
	}

	return result, nil
}

func createMap() map[rune]int {
	result := make(map[rune]int, len(shortLinkAlphabet))
	for i := 0; i < len(shortLinkAlphabet); i++ {
		result[shortLinkAlphabet[i]] = i
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
	id, err := stringToID(codedID)
	if err != nil {
		return "", err
	}

	fullLink, ok := s.lr.GetByID(id)
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return fullLink, nil
}
