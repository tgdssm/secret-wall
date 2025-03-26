package domain

import "time"

type User struct {
	ID          string    `json:"id"`
	AppleSub    string    `json:"apple_sub"`
	CreatedAt   time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
}

type UserRepo interface {
	FindByAppleSub(string) (*User, error)
	Create(*User) error
}
