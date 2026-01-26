package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerRunAddress *NetAddress
	ShortLinkAddress *NetAddress
}

type NetAddress struct {
	Protocol string
	Host     string
	Port     string
}

func (addr *NetAddress) String() string {
	return addr.Protocol + "://" + addr.Host + ":" + addr.Port + "/"
}

func (addr *NetAddress) Set(rawAddr string) error {
	if !strings.Contains(rawAddr, `://`) {
		rawAddr = `http://` + rawAddr
	}

	u, err := url.Parse(rawAddr)
	if err != nil {
		return errors.New("invalid address format")
	}

	protocol := u.Scheme
	host := u.Hostname()
	port := u.Port()

	if host == "" || port == "" {
		return errors.New("host or port is empty")
	}

	addr.Protocol = protocol
	addr.Host = host
	addr.Port = port
	return nil
}

var defaultServerRunAddress = &NetAddress{
	Protocol: "http",
	Host:     "localhost",
	Port:     "8080",
}

var defaultShortLinkAddress = &NetAddress{
	Protocol: "http",
	Host:     "localhost",
	Port:     "8080",
}

var cfg = &Config{
	ServerRunAddress: defaultServerRunAddress,
	ShortLinkAddress: defaultShortLinkAddress,
}

func GetConfig() *Config {
	if err := godotenv.Load(); err == nil {
		host, _ := os.LookupEnv("SHORTLINK_HTTP_URL")
		port, _ := os.LookupEnv("SHORTLINK_HTTP_PORT")
		cfg.ServerRunAddress.Host = host
		cfg.ServerRunAddress.Port = port
	}

	setConfigByFlags()
	return cfg
}

func setConfigByFlags() {
	flag.Var(cfg.ServerRunAddress, "a", "address to run the server on in format ip:port")
	flag.Var(cfg.ShortLinkAddress, "b", "default address for short links in format http(s)://ip:port/")
	flag.Parse()
}
