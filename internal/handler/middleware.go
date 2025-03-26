package handler

import (
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
