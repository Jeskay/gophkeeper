package abstraction

import "context"

type Repository interface {
	SetFilePath(value string)
	GetFilePath() string
}

type Client interface {
	StartUpload(ctx context.Context) (UploadStream, error)
}

type DataReader interface {
	ReadByChunk(name string, f func([]byte) error) error
}

type UploadStream interface {
	Init(fileName string, fileType string) error
	Upload(data []byte) error
	Close() (uint32, error)
}
