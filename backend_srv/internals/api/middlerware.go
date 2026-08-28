package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func authMiddlewareV1(next http.Handler) http.HandlerFunc {
	// return ServeHttp(http.ResponseWriter, *http.Request) SIGNATURE
	return func(w http.ResponseWriter, r *http.Request) {
		autorizationToken := r.Header.Get("authorization")

		log.Printf("logging authMiddleware: %s", autorizationToken)
		next.ServeHTTP(w, r)
	}
}

// only for httprouter
func authMiddlewareV2(next httprouter.Handle) httprouter.Handle {
	// return
	return func(w http.ResponseWriter, r *http.Request, path httprouter.Params) {
		autorizationToken := r.Header.Get("authorization")
		r.Header.Set("pass1", "pass1")

		log.Printf("logging authMiddleware: %s", autorizationToken)
		next(w, r, path)
	}
}

func level1(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, path httprouter.Params) {

		log.Printf("pass auth to <level1> : %s", r.Header.Get("pass1"))

		next(w, r, path)
	}
}
