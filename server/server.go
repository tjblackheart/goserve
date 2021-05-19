package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

type server struct {
	dir, port, certPath string
	caching, all, tls   bool
	signals             chan os.Signal
}

func New(dir, port string, caching, all, tls, force bool) (*server, error) {
	port, err := getPort(strings.Trim(port, ":"))
	if err != nil {
		return nil, err
	}

	dir, err = validateDir(dir)
	if err != nil {
		return nil, err
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := &server{
		dir:     dir,
		port:    port,
		caching: caching,
		all:     all,
		tls:     tls,
		signals: sigc,
	}

	if tls {
		userCache, err := os.UserCacheDir()
		if err != nil {
			return nil, errors.New("Could not get user cache: " + err.Error())
		}
		s.certPath = fmt.Sprintf("%s/goserve-certs", userCache)

		if err := generateCerts(s.certPath, force); err != nil {
			return nil, errors.New("Error generating certificates: " + err.Error())
		}
	}

	return s, nil
}

func (s server) Serve() {
	r := s.router()

	if s.all {
		log.Println("Binding to all interfaces ...")
	} else {
		log.Println("Binding to loopback interface ...")
	}

	ips, loopback, err := getIPs()
	if err != nil {
		log.Fatalln(err)
	}

	if loopback == "" && len(ips) == 0 {
		log.Fatalln("No interface to bind to")
	}

	if s.all {
		for _, ip := range ips {
			host := fmt.Sprintf("%s:%s", ip, s.port)
			srv := &http.Server{Handler: r, Addr: host}

			go func() {
				log.Printf("Serving from %s at http://%s ...", s.dir, host)
				log.Fatalln(srv.ListenAndServe())
			}()
		}
	}

	host := fmt.Sprintf("%s:%s", loopback, s.port)
	srv := &http.Server{Handler: r, Addr: host}
	proto := "http"
	if s.tls {
		proto = "https"
	}

	go func() {
		signal := <-s.signals
		fmt.Printf("%s received, shutting down\n", signal.String())
		srv.Shutdown(context.TODO())
	}()

	if s.tls {
		log.Printf("Serving secured from %s at %s://%s ...", s.dir, proto, host)
		key := s.certPath + "/server.key"
		cert := s.certPath + "/server.crt"
		log.Fatalln(srv.ListenAndServeTLS(cert, key))
	}

	log.Printf("Serving from %s at %s://%s ...", s.dir, proto, host)
	log.Fatalln(srv.ListenAndServe())
}

func (s server) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(recoverPanic, requestLog, gzip)

	if s.caching {
		r.Use(cache)
	}

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		path = filepath.Join(s.dir, path)
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(s.dir, "index.html"))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.FileServer(http.Dir(s.dir)).ServeHTTP(w, r)
	})

	return r
}
