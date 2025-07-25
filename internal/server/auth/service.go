package auth

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type Service interface {
	Login(context.Context, string, string) (int64, string)
	Register(context.Context, string, string) error
	VerifyToken(string) (*dto.User, error)
}
