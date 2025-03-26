package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
	ErrInvalidToken     = errors.New("invalid token")
)
