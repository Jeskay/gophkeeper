package file

import "path"

type fileService struct {
	fileWriter FileWriter
	fileReader FileReader
	prefix     string
}

func NewService(prefix string) *fileService {
	return &fileService{prefix: prefix}
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
