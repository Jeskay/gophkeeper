package control

import (
	"context"
	proto "gophkeeper/api/protos"
	menu "gophkeeper/internal/client/menu/control"
	"gophkeeper/internal/client/upload/abstraction"
	"path"
	"strings"
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

func (c *uploadController) UploadFile(filePath string) error {
	ctx := c.menuController.Authorize(context.Background())
	stream, err := c.client.StartUpload(ctx)
	if err != nil {
		return err
	}
	fileName := path.Base(filePath)
	str := strings.SplitN(fileName, ".", 2)
	stream.Init(str[0], "."+str[1])
	err = c.dataReader.ReadByChunk(filePath, func(data []byte) error {
		return stream.Upload(data)
	})
	_, respErr := stream.Close()
	if err == nil {
		return respErr
	}
	return err
}
