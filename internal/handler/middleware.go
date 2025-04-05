package handler

import (
	"log"
	"net/http"
	"secretWall/internal/domain"
	"secretWall/internal/service"
)

func Authenticator(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.ValidateToken(r); err != nil {
			WriteHTTPError(w, domain.ErrUserNotFound)
			return
		}
		handlerFunc(w, r)
	}
}

func Logger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s %s", r.Method, r.RequestURI, r.Host, r.Body)
		handlerFunc(w, r)
	}
}
