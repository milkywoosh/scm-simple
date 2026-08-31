package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"scm-simple-luke.com/dir/internals"
	"scm-simple-luke.com/dir/internals/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type ctxKey string

const payloadKey ctxKey = "payload"

func authMiddlewareV1(next http.Handler) http.HandlerFunc {
	// return ServeHttp(http.ResponseWriter, *http.Request) SIGNATURE
	return func(w http.ResponseWriter, r *http.Request) {
		autorizationToken := r.Header.Get("authorization")

		log.Printf("logging authMiddleware: %s", autorizationToken)
		next.ServeHTTP(w, r)
	}
}

// only for httprouter
func authMiddlewareV2(tokenMaker token.TokenMaker, next httprouter.Handle) httprouter.Handle {
	// return
	return func(w http.ResponseWriter, r *http.Request, path httprouter.Params) {
		autorizationToken := r.Header.Get(authorizationHeaderKey)
		// autorizationToken := "Bearer token "

		splittedWord := strings.Fields(autorizationToken)
		// for _, val := range splittedWord {
		// 	log.Printf(" ==> %s", val)
		// }

		errInfo := make(map[string]any)
		if len(splittedWord) < 2 {
			errInfo["message"] = "invalid token form, please re-check access token"
			internals.WriteErrorResponse(w, http.StatusUnauthorized, errInfo)
			return
		}

		payload, err := tokenMaker.VerifyToken(splittedWord[1])
		log.Printf("logging authMiddleware: %v", err)
		if err != nil {
			errInfo["message"] = err.Error()
			internals.WriteErrorResponse(w, http.StatusUnauthorized, errInfo)
			return
		}

		// passingCtx := context.WithValue(r.Context(), payloadKey, payload)

		// question, when will value inside this context EXPIRED??
		ctx := context.WithValue(r.Context(), payloadKey, payload)
		/*
			dataResp := make(map[string]any)
			dataResp["message"] = "sukses"
			dataResp["payload"] = payload
			internals.WriteResponse(w, http.StatusOK, dataResp)
		*/

		next(w, r.WithContext(ctx), path)
	}
}

func level1(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, path httprouter.Params) {

		log.Printf("pass auth to <level1> : %s", r.Header.Get("pass1"))

		next(w, r, path)
	}
}

func (s *Server) roleValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, path httprouter.Params) {

		log.Printf("roleValidationMiddleware : %s", r.Header.Get(authorizationPayloadKey))

		s.service.AuthenticationService.
		next(w, r, path)
	}
}
