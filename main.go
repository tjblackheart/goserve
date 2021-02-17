package main

import (
	"flag"
	"log"

	"github.com/tjblackheart/goserve/server"
)

func main() {
	var port string
	var caching bool
	var all bool
	var tls bool

	flag.StringVar(&port, "p", "9000", "Port")
	flag.BoolVar(&caching, "c", false, "Use http caching (default: false)")
	flag.BoolVar(&all, "a", false, "Bind to all interfaces (default: Loopback only)")
	flag.BoolVar(&tls, "s", false, "Use TLS. Will generate certs if they are not present (default: false)")
	flag.Parse()

	dir := flag.Arg(0)
	s, err := server.New(dir, port, caching, all, tls)
	if err != nil {
		log.Fatalln(err)
	}

	s.Serve()
}
