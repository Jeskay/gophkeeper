package auth

import "context"

type Service interface {
	Login(context.Context, string, string) (int64, string)
	Register(context.Context, string, string) error
	VerifyToken(string) (string, error)
}
