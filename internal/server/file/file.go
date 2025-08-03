package file

import (
	"gophkeeper/internal/server/dto"
	fPkg "gophkeeper/pkg/file"
	"io"
	"os"
	"path"
)

type fileService struct {
	fileWriter fPkg.FileWriter
	fileReader fPkg.FileReader
	prefix     string
}

const tmpSuffix = ".tmp"

func NewService(prefix string) *fileService {
	crntDir, err := os.Getwd()
	if err == nil {
		if _, err := os.Stat(prefix); os.IsNotExist(err) {
			if err = os.Mkdir(prefix, 0755); err != nil {
				prefix = crntDir
			}
		}
	}

	return &fileService{prefix: prefix, fileWriter: fPkg.NewFileWriter(), fileReader: fPkg.NewFileReader()}
}

func (s *fileService) SaveFile(name string, data []byte) error {
	file, err := s.fileWriter.CreateFile(path.Join(s.prefix, name))
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
	file, err := s.fileReader.OpenFile(path.Join(s.prefix, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []byte
	_, err = s.fileReader.FileRead(file, data)
	return data, err
}

func (s *fileService) WriteByChunk(name string, data <-chan dto.ChunkData) (size int, err error) {
	filePath := path.Join(s.prefix, name)
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
	var n int
	for chunk := range data {
		if chunk.Err != nil {
			err = chunk.Err
			return
		}
		n, err = s.fileWriter.BufferedWrite(writer, chunk.Data)
		if err != nil {
			return
		}
		size += n
	}
	s.fileWriter.BufferedFlush(writer)
	return
}

func (s *fileService) ReadByChunk(name string, data chan<- dto.ChunkData) error {
	filePath := path.Join(s.prefix, name)
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
