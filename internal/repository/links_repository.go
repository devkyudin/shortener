package repository

type LinksRepository struct {
	IdsLinks map[int]string
	LinksIds map[string]int
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		IdsLinks: make(map[int]string),
		LinksIds: make(map[string]int),
	}
}
func (lr *LinksRepository) CreateShortLink(originalLink string, id int) {
	lr.IdsLinks[id] = originalLink
	lr.LinksIds[originalLink] = id
}

func (lr *LinksRepository) GetByID(id int) (link string, isOk bool) {
	result, isOk := lr.IdsLinks[id]
	return result, isOk
}

func (lr *LinksRepository) GetByLink(link string) (id int, isOk bool) {
	result, isOk := lr.LinksIds[link]
	return result, isOk
}
