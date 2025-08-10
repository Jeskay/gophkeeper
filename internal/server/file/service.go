package file

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type Service interface {
	SaveFile(name string, data []byte) error
	ReadFile(name string) ([]byte, error)
	ReadByChunk(name string, data chan<- dto.ChunkData) error
	WriteByChunk(ctx context.Context, name string, data <-chan dto.ChunkData) (size int, err error)
}
