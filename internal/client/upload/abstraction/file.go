package abstraction

import (
	"gophkeeper/pkg/file"
)

type dataReader struct {
	reader file.FileReader
}

func NewDataReader() *dataReader {
	return &dataReader{reader: file.NewFileReader()}
}

func (r *dataReader) ReadByChunk(name string, f func([]byte) error) error {
	file, err := r.reader.OpenFile(name)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := r.reader.NewBufferedScanner(file)
	for scanner.Scan() {
		if err := f(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
