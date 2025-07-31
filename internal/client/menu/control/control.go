package control

import "context"

type Controller interface {
	IsAuthenticated() bool
	SaveToken(value string)
	Authorize(ctx context.Context) context.Context
}
