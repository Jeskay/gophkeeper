package control

import (
	"context"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/client/download/abstraction"
)

type downloadController struct {
	client     abstraction.Client
	repository abstraction.Repository
}

func NewController(grpcClient proto.GophkeeperClient) *downloadController {
	return &downloadController{
		client:     abstraction.NewClient(grpcClient),
		repository: abstraction.NewRepository(),
	}
}

func (c *downloadController) DownloadFiles(ctx context.Context) ([]abstraction.FileData, error) {
	files, err := c.client.GetFiles(ctx)
	if err != nil {
		return nil, err
	}
	c.repository.SetFiles(files)
	return files, nil
}
