package model

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type CodedLink struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewCodedLink(id int, originalURL string, alphabet *Alphabet) *CodedLink {
	var result = ""
	originalID := id
	alphabetLength := len(alphabet.Chars)
	for id > 0 {
		code := id % alphabetLength
		result = result + string(alphabet.Chars[code])
		id = id / alphabetLength
	}

	shortURL := reverse(result)

	return &CodedLink{
		UUID:        originalID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func FromJSON(link string) (*CodedLink, error) {
	result := CodedLink{}
	err := json.Unmarshal([]byte(link), &result)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить строку в структуру CodedLink: %w", err)
	}

	return &result, nil
}

func (codedLink *CodedLink) ToJSON() ([]byte, error) {
	return json.Marshal(codedLink)
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
