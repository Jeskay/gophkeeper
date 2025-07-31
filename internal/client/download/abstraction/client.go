package abstraction

import (
	"context"
	proto "gophkeeper/api/protos"
)

type downloadClient struct {
	grpcClient proto.GophkeeperClient
}

func NewClient(grpcClient proto.GophkeeperClient) *downloadClient {
	return &downloadClient{grpcClient: grpcClient}
}

func (c *downloadClient) GetFiles(ctx context.Context) ([]FileData, error) {
	return nil, nil
}
