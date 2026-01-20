package testutils

import (
	"os"

	"github.com/devkyudin/shortener/internal/config"
)

func SetupTestEnvironment() {
	os.Args = append(os.Args, "-a", "localhost:8080", "-b", "localhost:8080/")
	config.ParseFlags()
}
