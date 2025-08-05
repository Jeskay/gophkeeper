package keeper

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type Service interface {
	GetFiles(ctx context.Context, userId int64) ([]dto.File, error)
	UploadFile(ctx context.Context, userId int64, file dto.File, chData <-chan dto.ChunkData) (int, error)
	DownloadFile(ctx context.Context, userId int64, file dto.File, out chan<- dto.ChunkData) (*dto.File, error)
}
