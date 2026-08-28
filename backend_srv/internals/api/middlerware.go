package api

import (
	"log"
	"net/http"
)

func authMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autorizationToken := r.Header.Get("authorization")

		log.Printf("logging authMiddleware: %s", autorizationToken)
		next.ServeHTTP(w, r)
	}
}
