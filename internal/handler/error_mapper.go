package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"secretWall/internal/domain"
)

func MapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusUnauthorized, "user not found"
	case errors.Is(err, domain.ErrUserAlreadyExist):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, domain.ErrInvalidToken):
		return http.StatusBadRequest, "invalid token"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	status, msg := MapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
