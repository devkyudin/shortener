package main

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.ParseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(config.FlagRunAddress, Router())
}

func Router() chi.Router {
	r := chi.NewRouter()
	r.HandleFunc(`/`, handler.Shorten)
	r.HandleFunc(`/{id}`, handler.GetLink)
	return r
}
