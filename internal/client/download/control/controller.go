package control

import (
	"bytes"
	"context"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/client/download/abstraction"
	menu "gophkeeper/internal/client/menu/control"
	"io"
)

type downloadController struct {
	menuController menu.Controller
	client         abstraction.Client
	repository     abstraction.Repository
	dataWriter     abstraction.DataWriter
}

func NewController(saveDir string, grpcClient proto.GophkeeperClient, menuController menu.Controller) *downloadController {
	return &downloadController{
		client:         abstraction.NewClient(grpcClient),
		repository:     abstraction.NewRepository(),
		menuController: menuController,
		dataWriter:     abstraction.NewDataWriter(saveDir),
	}
}

func (c *downloadController) DownloadFileList() ([]*abstraction.FileData, error) {
	ctx := c.menuController.Authorize(context.Background())
	files, err := c.client.GetFiles(ctx)
	if err != nil {
		return nil, err
	}
	c.repository.SetFiles(files)
	return files, nil
}

// TODO: make download by file ID
func (c *downloadController) DownloadFile(id int64) error {
	ctx := c.menuController.Authorize(context.Background())
	stream, err := c.client.StartDownload(ctx, id)
	if err != nil {
		return err
	}
	fileName, err := stream.Init()
	if err != nil {
		return err
	}
	var data bytes.Buffer
	for {
		chunk, err := stream.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = data.Write(chunk)
		if err != nil {
			return err
		}
	}
	return c.dataWriter.WriteFile(fileName, data.Bytes())
}
