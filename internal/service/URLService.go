package service

import (
	"errors"
	"sync"
	"unicode/utf8"
)

var m = make(map[int]string)
var newShortLinkID = initialShortLinkID
var hostAddress = `http://localhost:8080/`
var initialShortLinkID = 10_000_000
var mutex sync.Mutex
var shortLinkAlphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var alphabetMap map[rune]int

func CreateShortLink(originalLink string) string {
	id := GetNewID()
	m[id] = originalLink
	shortedLink := hostAddress + IDToString(id)
	return shortedLink
}

func GetNewID() int {
	mutex.Lock()
	defer mutex.Unlock()
	result := newShortLinkID
	newShortLinkID++
	return result
}

func IDToString(id int) string {
	var result = ""
	alphabetLength := len(shortLinkAlphabet)
	for id > 0 {
		code := id % alphabetLength
		result = result + string(shortLinkAlphabet[code])
		id = id / alphabetLength
	}

	return Reverse(result)
}

func StringToID(src string) (int, error) {
	if alphabetMap == nil {
		alphabetMap = CreateMap()
	}
	result := 0
	alphabetLength := len(shortLinkAlphabet)
	reversed := Reverse(src)
	runes := []rune(reversed)
	for i := 0; i < len(src); i++ {
		runeID, ok := alphabetMap[runes[i]]
		if !ok {
			return 0, errors.New(`битая ссылка`)
		}
		result += runeID * IntPow(alphabetLength, i)
	}

	return result, nil
}

func CreateMap() map[rune]int {
	result := make(map[rune]int, len(shortLinkAlphabet))
	for i := 0; i < len(shortLinkAlphabet); i++ {
		result[shortLinkAlphabet[i]] = i
	}

	return result
}

func IntPow(n, m int) int {
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

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func GetFullLink(codedID string) (string, error) {
	id, err := StringToID(codedID)
	if err != nil {
		return "", err
	}

	fullLink, ok := m[id]
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return fullLink, nil
}
