package abstraction

import "context"

type Client interface {
	Register(ctx context.Context, name, password string) error
	Authenticate(ctx context.Context, name, password string) (string, error)
}
