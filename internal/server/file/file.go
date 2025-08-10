package file

import (
	"context"
	"io"
	"os"
	"path"

	"gophkeeper/internal/server/dto"

	fPkg "gophkeeper/pkg/file"
)

type fileService struct {
	fileWriter fPkg.FileWriter
	fileReader fPkg.FileReader
	storePath     string
}

const tmpSuffix = ".tmp"

func NewService(storePath string) *fileService {
		return &fileService{storePath: storePath, fileWriter: fPkg.NewFileWriter(), fileReader: fPkg.NewFileReader()}
}

func (s *fileService) SaveFile(name string, data []byte) error {
	file, err := s.fileWriter.CreateFile(path.Join(s.storePath, name))
	if err != nil {
		return err
	}
	defer s.fileWriter.FileClose(file)

	_, err = s.fileWriter.FileWrite(file, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *fileService) ReadFile(name string) ([]byte, error) {
	file, err := s.fileReader.OpenFile(path.Join(s.storePath, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []byte
	_, err = s.fileReader.FileRead(file, data)
	return data, err
}

func (s *fileService) WriteByChunk(ctx context.Context, name string, data <-chan dto.ChunkData) (size int, err error) {
	filePath := path.Join(s.storePath, name)
	file, err := os.Create(filePath + tmpSuffix)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
		if err == nil {
			err = os.Rename(file.Name(), filePath)
		}
	}()
	writer := s.fileWriter.NewBufferedWriter(file)
	defer s.fileWriter.BufferedFlush(writer)
	var n int
	for {
		select {
		case <-ctx.Done():
			err = context.Canceled
			return
		case chunk, ok := <-data:
			if !ok || chunk.Err != nil {
				err = chunk.Err
				return
			}
			n, err = s.fileWriter.BufferedWrite(writer, chunk.Data)
			if err != nil {
				return
			}
			size += n
		}
	}
}

func (s *fileService) ReadByChunk(name string, data chan<- dto.ChunkData) error {
	filePath := path.Join(s.storePath, name)
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	file, err := s.fileReader.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := s.fileReader.NewBufferedReader(file)
	for {
		chunk, err := s.fileReader.BufferedRead(reader)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		data <- dto.ChunkData{Data: chunk}
	}
	return nil
}
