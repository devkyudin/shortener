package main

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/handler"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.Shorten)
	mux.HandleFunc(`/{id}`, handler.GetLink)
	return http.ListenAndServe(`:8080`, mux)
}
