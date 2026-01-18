package config

import "flag"

var FlagRunAddress string
var FlagDefaultAddress string

func ParseFlags() {
	flag.StringVar(&FlagRunAddress, "a", ":8080", "address to run the server on")
	flag.StringVar(&FlagDefaultAddress, "b", "http://localhost:8080/", "default address for short links")
	flag.Parse()
}
