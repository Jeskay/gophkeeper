package control

import (
	"gophkeeper/pkg/file"
	"path"
)

type dataWriter struct {
	writer   file.FileWriter
	savePath string
}

func NewDataWriter(savePath string) *dataWriter {
	return &dataWriter{writer: file.NewFileWriter(), savePath: savePath}
}

func (w *dataWriter) WriteFile(name string, data []byte) error {
	file, err := w.writer.CreateFile(path.Join(w.savePath, name))
	if err != nil {
		return err
	}
	defer w.writer.FileClose(file)

	_, err = w.writer.FileWrite(file, data)
	if err != nil {
		return err
	}
	return nil
}
