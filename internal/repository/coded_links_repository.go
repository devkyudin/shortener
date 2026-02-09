package repository

import (
	"sync"

	"github.com/devkyudin/shortener/internal/model"
)

type CodedLinksRepository struct {
	idsMap          map[int]model.CodedLink
	shortUrlsMap    map[string]model.CodedLink
	originalUrlsMap map[string]model.CodedLink
	nextID          int
	rwMutex         *sync.RWMutex
}

func NewLinksRepository() *CodedLinksRepository {
	return &CodedLinksRepository{
		idsMap:          make(map[int]model.CodedLink),
		shortUrlsMap:    make(map[string]model.CodedLink),
		originalUrlsMap: make(map[string]model.CodedLink),
		nextID:          1_000_000,
		rwMutex:         &sync.RWMutex{},
	}
}
func (lr *CodedLinksRepository) CreateCodedLink(codedLink model.CodedLink) model.CodedLink {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	if found, isOk := lr.shortUrlsMap[codedLink.OriginalUrl]; isOk {
		return found
	}
	codedLink.UUID = lr.nextID
	lr.idsMap[codedLink.UUID] = codedLink
	lr.shortUrlsMap[codedLink.ShortUrl] = codedLink
	lr.originalUrlsMap[codedLink.OriginalUrl] = codedLink
	return codedLink
}

func (lr *CodedLinksRepository) GetByID(id int) (link model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.idsMap[id]
	return result, isOk
}

func (lr *CodedLinksRepository) GetByShortUrl(shortLink string) (link model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.shortUrlsMap[shortLink]
	return result, isOk
}

func (lr *CodedLinksRepository) GetByOriginalUrl(originalLink string) (link model.CodedLink, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.originalUrlsMap[originalLink]
	return result, isOk
}

func (lr *CodedLinksRepository) GetUniqueID() int {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	result := lr.nextID
	lr.nextID++
	return result
}
