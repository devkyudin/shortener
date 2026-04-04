package repository

import "github.com/devkyudin/shortener/internal/model"

type CodedLinksRepository interface {
	CreateCodedLink(codedLink *model.CodedLink) error
	GetByID(id int) (link *model.CodedLink, isOk bool)
	GetByShortURL(shortLink string) (link *model.CodedLink, isOk bool)
	GetByOriginalURL(originalLink string) (link *model.CodedLink, isOk bool)
	GetUniqueID() int
}
