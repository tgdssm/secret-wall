package service

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"io"
	"math/big"
	"net/http"
	"secretWall/internal/domain"
	"time"
)

type AuthService struct {
	userRepo domain.UserRepo
}

type ApplePublicKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type AppleKeysResponse struct {
	Keys []ApplePublicKey `json:"keys"`
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

		var keys AppleKeysResponse
		err = json.Unmarshal(body, &keys)
		if err != nil {
			return nil, err
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, domain.ErrInvalidToken
		}

		for _, key := range keys.Keys {
			if key.Kid == kid {
				nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, err
				}
				eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
				if err != nil {
					return nil, err
				}

				e := 0
				for _, b := range eBytes {
					e = e<<8 + int(b)
				}

				publicKey := &rsa.PublicKey{
					N: new(big.Int).SetBytes(nBytes),
					E: e,
				}

				return publicKey, nil
			}
		}

		return nil, domain.ErrInvalidToken
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.userRepo.FindByAppleSub(sub)

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
