package service

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"secretWall/internal/domain"
	"time"
)

type AuthService struct {
	userRepo domain.UserRepo
}

func NewAuthService(userRepo domain.UserRepo) AuthService {
	return AuthService{userRepo: userRepo}
}

func (s *AuthService) AuthenticateWithApple(identityToken string) (*domain.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(identityToken, claims, func(token *jwt.Token) (interface{}, error) {
		resp, err := http.Get("https://appleid.apple.com/auth/keys")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var keys map[string]interface{}
		err = json.Unmarshal(body, &keys)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.FindByAppleSub(sub)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &domain.User{
			ID:        generateUUID(),
			AppleSub:  sub,
			CreatedAt: time.Now(),
		}

		if err := s.userRepo.Create(user); err != nil {
			if errors.Is(err, domain.ErrUserAlreadyExist) {
				return nil, err
			}
			return nil, err
		}
	}

	accessToken, err := CreateToken(user.ID)
	user.AccessToken = accessToken

	return user, nil
}

func generateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}
