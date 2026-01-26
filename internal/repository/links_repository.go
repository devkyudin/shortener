package repository

import "sync"

type LinksRepository struct {
	IdsLinks map[int]string
	LinksIds map[string]int
	rwMutex  *sync.RWMutex
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		IdsLinks: make(map[int]string),
		LinksIds: make(map[string]int),
		rwMutex:  &sync.RWMutex{},
	}
}
func (lr *LinksRepository) CreateShortLink(originalLink string, id int) {
	lr.rwMutex.Lock()
	defer lr.rwMutex.Unlock()
	lr.IdsLinks[id] = originalLink
	lr.LinksIds[originalLink] = id
}

func (lr *LinksRepository) GetByID(id int) (link string, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.IdsLinks[id]
	return result, isOk
}

func (lr *LinksRepository) GetByLink(link string) (id int, isOk bool) {
	lr.rwMutex.RLock()
	defer lr.rwMutex.RUnlock()
	result, isOk := lr.LinksIds[link]
	return result, isOk
}
