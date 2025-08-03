package abstraction

import "context"

type Repository interface {
	SetFilePath(value string)
	GetFilePath() string
}

type Client interface {
	StartUpload(ctx context.Context) (UploadStream, error)
}

type UploadStream interface {
	Init(fileName string, fileType string, size uint32) error
	Upload(data []byte) error
	Close() (uint32, error)
}
