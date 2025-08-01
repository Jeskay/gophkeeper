package file

import (
	fPkg "gophkeeper/pkg/file"
	"os"
	"path"
)

type fileService struct {
	fileWriter fPkg.FileWriter
	fileReader fPkg.FileReader
	prefix     string
}

func NewService(prefix string) *fileService {
	crntDir, err := os.Getwd()
	if err == nil {
		prefix = crntDir //TODO: check for directory and store in prefix folder
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

func (s *fileService) ReadByChunk(name string, f func([]byte) error) error {
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

	scanner := s.fileReader.NewBufferedScanner(file)
	for scanner.Scan() {
		if err := f(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
