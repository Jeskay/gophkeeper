package control

import (
	"context"
	pb "gophkeeper/api/protos"
	"gophkeeper/internal/client/authentication/abstraction"

	"google.golang.org/grpc/metadata"
)

type authController struct {
	client abstraction.Client
}

func NewController(grpcClient pb.GophkeeperClient) *authController {
	return &authController{
		client: abstraction.NewClient(grpcClient),
	}
}

func (c *authController) Register(name, password string) error {
	return c.client.Register(context.TODO(), name, password)
}

func (c *authController) Authenticate(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "jwt", token)
}

func (c *authController) LogIn(name, password string) (string, error) {
	return c.client.Authenticate(context.TODO(), name, password)
}
