package file

import (
	"bufio"
	"os"
)

type fileReader struct{}

func NewFileReader() FileReader {
	return &fileReader{}
}

func (r *fileReader) OpenFile(name string) (*os.File, error) {
	return os.Open(name)
}

func (r *fileReader) FileClose(file *os.File) error {
	return file.Close()
}

func (r *fileReader) FileRead(file *os.File, b []byte) (n int, err error) {
	return file.Read(b)
}

func (r *fileReader) NewBufferedScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 2048)
	return scanner
}

func (r *fileReader) NewBufferedReader(file *os.File) *bufio.Reader {
	return bufio.NewReader(file)
}

func (r *fileReader) BufferedRead(reader *bufio.Reader) ([]byte, error) {
	data := make([]byte, 100)
	n, err := reader.Read(data)
	return data[:n], err
}
func (r *fileReader) ReadByChunk(name string, f func([]byte) error) error {
	file, err := r.OpenFile(name)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := r.NewBufferedScanner(file)
	for scanner.Scan() {
		if err := f(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
