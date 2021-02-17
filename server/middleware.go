package server

import (
	"fmt"
	"log"
	"net/http"

	gz "github.com/NYTimes/gziphandler"
)

func gzip(next http.Handler) http.Handler {
	return gz.GzipHandler(next)
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "max-age=2592000")
		next.ServeHTTP(w, r)
	})
}

func recoverErr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
