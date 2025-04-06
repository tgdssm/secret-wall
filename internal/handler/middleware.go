package handler

import (
	"io"
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
		log.Printf("\n %s %s", r.Method, r.RequestURI)
		body, _ := io.ReadAll(r.Body)
		log.Printf("\n %s", body)
		handlerFunc(w, r)
	}
}
