package control

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strings"

	proto "gophkeeper/api/protos"
	menu "gophkeeper/internal/client/menu/control"
	"gophkeeper/internal/client/upload/abstraction"
	"gophkeeper/pkg/cipher"
	"gophkeeper/pkg/file"
)

type uploadController struct {
	menuController menu.Controller
	client         abstraction.Client
	repository     abstraction.Repository
	fileReader     file.FileReader
	fileWriter     file.FileWriter
	uploadCipher   cipher.Cipher
}

const encryptedSuffix = ".tmp"

func NewController(grpcClient proto.GophkeeperClient, uploadCipher cipher.Cipher, menuController menu.Controller) *uploadController {

	return &uploadController{
		menuController: menuController,
		client:         NewClient(grpcClient),
		repository:     abstraction.NewRepository(),
		fileReader:     file.NewFileReader(),
		uploadCipher:   uploadCipher,
		fileWriter:     file.NewFileWriter(),
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
	stream.Init(strings.TrimSuffix(fileName, fileExt), fileExt, uint32(info.Size()))
	encrypted, err := c.encryptToFile(filePath)
	if err != nil {
		return err
	}
	f, err := c.fileReader.OpenFile(encrypted)
	if err != nil {
		return err
	}
	reader := c.fileReader.NewBufferedReader(f)
	for {
		b, err := c.fileReader.BufferedRead(reader)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = stream.Upload(b)
		if err != nil {
			return err
		}
	}
	//defer os.Remove(encrypted)
	_, err = stream.Close()
	return err
}

func (c *uploadController) encryptToFile(fileName string) (out string, err error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return
	}
	var ciphered []byte
	ciphered, err = c.uploadCipher.Encrypt(data)
	if err != nil {
		return
	}
	out = fmt.Sprint(fileName, encryptedSuffix)
	var cFile *os.File
	cFile, err = c.fileWriter.CreateFile(out)
	if err != nil {
		return
	}
	defer cFile.Close()
	_, err = c.fileWriter.FileWrite(cFile, ciphered)
	if err != nil {
		return "", err
	}
	return
}
