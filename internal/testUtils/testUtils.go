package testUtils

import (
	"os"

	"github.com/devkyudin/shortener/internal/config"
)

//func SetupTestEnvironment() {
//	argsMap := make(map[string]struct{})
//	for i := range os.Args {
//		argsMap[os.Args[i]] = struct{}{}
//	}
//
//
//	argsMap["-a"] = struct{}{}
//	argsMap[":8080"] = struct{}{}
//	argsMap["-b"] = struct{}{}
//	argsMap["http://localhost:8080/"] = struct{}{}
//	os.Args = make([]string, len(argsMap))
//	i := 0
//	for v := range maps.Keys(argsMap) {
//		os.Args[i] = v
//		i++
//	}
//	config.ParseFlags()
//}

func SetupTestEnvironment() {
	os.Args = append(os.Args, "-a", ":8080", "-b", "http://localhost:8080/")
	config.ParseFlags()
}
