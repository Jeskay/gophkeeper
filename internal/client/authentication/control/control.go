package control

import "context"

type Controller interface {
	Register(name, password string) error
	Authenticate(ctx context.Context, token string) context.Context
	LogIn(name, password string) (string, error)
}
