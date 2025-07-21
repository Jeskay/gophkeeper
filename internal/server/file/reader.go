package file

import (
	"os"
)

type FileReader interface {
	OpenFile(name string) (*os.File, error)
	FileRead(file *os.File, b []byte) (n int, err error)
}

type fileReader struct{}

func NewFileReader() FileReader {
	return &fileReader{}
}

func (r *fileReader) OpenFile(name string) (*os.File, error) {
	return os.Open(name)
}

func (r *fileReader) FileRead(file *os.File, b []byte) (n int, err error) {
	return file.Read(b)
}
