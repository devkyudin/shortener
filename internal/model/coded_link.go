package model

import "unicode/utf8"

type CodedLink struct {
	UUID        int    `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func NewCodedLink(id int, originalUrl string, alphabet *Alphabet) CodedLink {
	var result = ""
	alphabetLength := len(alphabet.Chars)
	for id > 0 {
		code := id % alphabetLength
		result = result + string(alphabet.Chars[code])
		id = id / alphabetLength
	}

	shortUrl := reverse(result)

	return CodedLink{
		UUID:        id,
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
	}
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
