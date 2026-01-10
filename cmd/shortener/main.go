package main

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"unicode/utf8"
)

var m map[int]string = make(map[int]string)
var hostAddress string = `localhost:8080/`
var initialShortLinkID = 10_000_000
var newShortLinkID = initialShortLinkID
var mutex sync.Mutex
var shortLinkAlphabet []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var alphabetMap map[rune]int

func main() {
	alphabetMap = CreateMap()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, Shorten)
	mux.HandleFunc(`/{id}`, GetLink)
	return http.ListenAndServe(`:8080`, mux)
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get(`Content-Type`) != `text/plain` {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalLink, err := io.ReadAll(r.Body)
	if err != nil || string(originalLink) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortedLink := CreateShortLink(string(originalLink))
	w.Header().Add(`Content-Type`, `text/plain`)
	w.Header().Add(`Location`, shortedLink)
	w.WriteHeader(http.StatusCreated)
}

func GetLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.PathValue(`id`)
	result, error := GetFullLink(id)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add(`Content-Type`, `text/plain`)
	w.Header().Add(`Location`, result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func CreateShortLink(originalLink string) string {
	id := GetNewID()
	m[id] = originalLink
	shortedLink := hostAddress + IDToString(id)
	return shortedLink
}

func GetFullLink(codedID string) (string, error) {
	id, error := StringToID(codedID)
	if error != nil {
		return "", error
	}

	fullLink, ok := m[id]
	if !ok {
		return "", errors.New(`нет ссылки с таким идентификатором идентификатором`)
	}

	return fullLink, nil
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

// IntPow calculates n to the mth power. Since the result is an int, it is assumed that m is a positive power
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

func CreateMap() map[rune]int {
	result := make(map[rune]int, len(shortLinkAlphabet))
	for i := 0; i < len(shortLinkAlphabet); i++ {
		result[shortLinkAlphabet[i]] = i
	}

	return result
}

func GetNewID() int {
	mutex.Lock()
	defer mutex.Unlock()
	result := newShortLinkID
	newShortLinkID++
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
