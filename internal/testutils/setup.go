package testutils

import (
	"os"

	"github.com/devkyudin/shortener/internal/config"
)

func SetupTestEnvironment() {
	os.Args = append(os.Args, "-a", ":8080", "-b", "http://localhost:8080/")
	config.ParseFlags()
}
