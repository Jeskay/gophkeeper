package auth

import (
	"context"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/dto"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	dbService db.Service
	secretKey []byte
	expiresAt time.Duration
}

func NewService(dbService db.Service, secretKey []byte, expiration time.Duration) *authService {
	return &authService{
		secretKey: secretKey,
		expiresAt: expiration,
		dbService: dbService,
	}
}

func (s *authService) Login(ctx context.Context, name, password string) (int64, string) {
	user, err := s.dbService.GetUser(ctx, name)
	if err != nil {
		return http.StatusNotFound, ""
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return http.StatusUnauthorized, ""
	}
	token, err := s.CreateToken(user.Name, user.Id)
	if err != nil {
		return http.StatusUnauthorized, ""
	}
	return http.StatusAccepted, token
}

func (s *authService) Register(ctx context.Context, name, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}
	err = s.dbService.CreateUser(ctx, dto.User{Name: name, Password: string(hash)})
	return err
}
