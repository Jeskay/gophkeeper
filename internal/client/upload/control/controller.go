package control

import (
	"context"
	proto "gophkeeper/api/protos"
	menu "gophkeeper/internal/client/menu/control"
	"gophkeeper/internal/client/upload/abstraction"
)

type uploadController struct {
	menuController menu.Controller
	client         abstraction.Client
	repository     abstraction.Repository
	dataReader     abstraction.DataReader
}

func NewController(grpcClient proto.GophkeeperClient, menuController menu.Controller) *uploadController {
	return &uploadController{
		menuController: menuController,
		client:         abstraction.NewClient(grpcClient),
		repository:     abstraction.NewRepository(),
		dataReader:     abstraction.NewDataReader(),
	}
}

func (c *uploadController) UploadFile(path string) error {
	ctx := c.menuController.Authorize(context.Background())
	stream, err := c.client.StartUpload(ctx)
	if err != nil {
		return err
	}
	stream.Init("test1", ".txt")
	err = c.dataReader.ReadByChunk(path, func(data []byte) error {
		return stream.Upload(data)
	})
	_, respErr := stream.Close()
	if err == nil {
		return respErr
	}
	return err
}
