package file

import (
	"bufio"
	"os"
)

type FileReader interface {
	OpenFile(name string) (*os.File, error)
	FileRead(file *os.File, b []byte) (n int, err error)
	NewBufferedScanner(file *os.File) *bufio.Scanner
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

func (r *fileReader) NewBufferedScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 2048)
	return scanner
}
