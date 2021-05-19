package main

import (
	"flag"
	"log"

	"github.com/tjblackheart/goserve/server"
)

func main() {
	var (
		port    string
		caching bool
		all     bool
		tls     bool
		force   bool
	)

	flag.StringVar(&port, "p", "9000", "Port")
	flag.BoolVar(&caching, "c", false, "Use http caching")
	flag.BoolVar(&all, "a", false, "Bind to all interfaces (default: Loopback only)")
	flag.BoolVar(&tls, "s", false, "Use TLS. Will generate certs if they are not present")
	flag.BoolVar(&force, "f", false, "Force certificate generation")
	flag.Parse()

	dir := flag.Arg(0)
	s, err := server.New(dir, port, caching, all, tls, force)
	if err != nil {
		log.Fatalln(err)
	}

	s.Serve()
}
