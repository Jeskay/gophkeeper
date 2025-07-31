package abstraction

import (
	"context"
	proto "gophkeeper/api/protos"
)

type authClient struct {
	grpcClient proto.GophkeeperClient
}

func NewClient(grpcClient proto.GophkeeperClient) *authClient {
	return &authClient{grpcClient: grpcClient}
}

func (c *authClient) Register(ctx context.Context, name, password string) error {
	_, err := c.grpcClient.Register(ctx, &proto.RegisterRequest{Name: name, Password: password})
	if err != nil {
		return err
	}
	return nil
}

func (c *authClient) Authenticate(ctx context.Context, name, password string) (string, error) {
	res, err := c.grpcClient.Login(ctx, &proto.LoginRequest{Name: name, Password: password})
	if err != nil {
		return "", err
	}
	return res.Token, nil
}
