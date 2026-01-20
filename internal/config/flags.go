package config

import (
	"errors"
	"flag"
	"net/url"
	"strings"
)

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

var Cfg = Config{
	ServerRunAddress:   &NetAddress{},
	FlagDefaultAddress: &NetAddress{},
}

type Config struct {
	ServerRunAddress   *NetAddress
	FlagDefaultAddress *NetAddress
}

func ParseFlags() {
	flag.Var(Cfg.ServerRunAddress, "a", "address to run the server on in format ip:port")
	flag.Var(Cfg.FlagDefaultAddress, "b", "default address for short links in format http(s)://ip:port/")
	flag.Parse()
}
