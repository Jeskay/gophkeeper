package abstraction

import "context"

type Repository interface {
	SetFiles(value []*FileData)
	GetFiles() []*FileData
}

type DataWriter interface {
	WriteFile(name string, data []byte) error
}

type Client interface {
	GetFiles(ctx context.Context) ([]*FileData, error)
	StartDownload(ctx context.Context, id int64) (DownloadStream, error)
}

type DownloadStream interface {
	Init() (string, error)
	Receive() ([]byte, error)
	Close() error //TODO: delete since unnecessary when reading stream
}
