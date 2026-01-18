package repository

var idsLinks = make(map[int]string)
var linksIds = make(map[string]int)

func CreateShortLink(originalLink string, id int) {
	idsLinks[id] = originalLink
	linksIds[originalLink] = id
}

func GetById(id int) (link string, isOk bool) {
	result, isOk := idsLinks[id]
	return result, isOk
}

func GetByLink(link string) (id int, isOk bool) {
	result, isOk := linksIds[link]
	return result, isOk
}
