package control

import (
	"context"
	"errors"
	proto "gophkeeper/api/protos"
	menu "gophkeeper/internal/client/menu/control"
	"gophkeeper/internal/client/upload/abstraction"
	"math"
	"os"
	"path"
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
		client:         NewClient(grpcClient),
		repository:     abstraction.NewRepository(),
		dataReader:     NewDataReader(),
	}
}

func (c *uploadController) UploadFile(filePath string) error {
	ctx := c.menuController.Authorize(context.Background())
	stream, err := c.client.StartUpload(ctx)
	if err != nil {
		return err
	}
	fileName := path.Base(filePath)
	fileExt := path.Ext(filePath)
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.Size() > math.MaxUint32 {
		return errors.New("file is too large")
	}
	stream.Init(fileName, fileExt, uint32(info.Size()))
	err = c.dataReader.ReadByChunk(filePath, func(data []byte) error {
		return stream.Upload(data)
	})
	_, respErr := stream.Close()
	if err == nil {
		return respErr
	}
	return err
}
