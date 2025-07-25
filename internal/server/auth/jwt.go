package auth

import (
	"errors"
	"gophkeeper/internal/server/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Login  string
	UserId int64
}

func (s *authService) CreateToken(login string, userId int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresAt)),
		},
		Login:  login,
		UserId: userId,
	})

	return claims.SignedString(s.secretKey)
}

func (s *authService) VerifyToken(tokenString string) (*dto.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &dto.User{Name: claims.Login, Id: claims.UserId}, nil
}
