package abstraction

import "context"

type Repository interface {
	SetFiles(value []FileData)
	GetFiles() []FileData
}

type Client interface {
	GetFiles(ctx context.Context) ([]FileData, error)
}
