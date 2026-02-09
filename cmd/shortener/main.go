package main

import (
	"net/http"

	"github.com/devkyudin/shortener/internal/dependencies"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	deps := dependencies.GetDependencies()
	serverAddress := deps.Config.ServerRunAddress.Host + ":" + deps.Config.ServerRunAddress.Port
	return http.ListenAndServe(serverAddress, *deps.Router)
}
