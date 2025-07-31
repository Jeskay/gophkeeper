package control

import (
	"context"
	"gophkeeper/internal/client/download/abstraction"
)

type Controller interface {
	DownloadFiles(ctx context.Context) ([]abstraction.FileData, error)
}
