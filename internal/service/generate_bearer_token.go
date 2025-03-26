package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"secretWall/internal/domain"
	"strings"
	"time"
)

func CreateToken(userId string) (string, error) {
	permission := jwt.MapClaims{}
	permission["sub"] = userId
	permission["authorized"] = true
	permission["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permission)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, getVerificationKey)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return domain.ErrInvalidToken
}

func getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	splitToken := strings.Split(bearToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
