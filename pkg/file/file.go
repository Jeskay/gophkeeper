package file

import (
	"bufio"
	"io/fs"
	"os"
)

type FileReader interface {
	OpenFile(name string) (*os.File, error)
	FileRead(file *os.File, b []byte) (n int, err error)
	FileClose(file *os.File) error
	NewBufferedReader(file *os.File) *bufio.Reader
	BufferedRead(reader *bufio.Reader) ([]byte, error)
	NewBufferedScanner(file *os.File) *bufio.Scanner
	ReadByChunk(name string, f func([]byte) error) error
}

type FileWriter interface {
	CreateFile(name string) (*os.File, error)
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	FileWriteString(file *os.File, s string) (n int, err error)
	FileWrite(file *os.File, b []byte) (n int, err error)
	FileClose(file *os.File) error
	NewBufferedWriter(file *os.File) *bufio.Writer
	BufferedWriteString(writer *bufio.Writer, s string) (n int, err error)
	BufferedFlush(writer *bufio.Writer) error
}
