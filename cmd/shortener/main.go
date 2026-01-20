package main

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/handler"
)

func main() {
	config.ParseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(config.FlagRunAddress, handler.ShortenerRouter())
}
