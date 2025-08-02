package db

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type Service interface {
	GetUser(ctx context.Context, name string) (dto.User, error)
	CreateUser(ctx context.Context, user dto.User) error
	GetFile(ctx context.Context, userId int64, fileId int64) (dto.File, error)
	CreateFile(ctx context.Context, file dto.File, userId int64) (int64, error)
	GetUserFiles(ctx context.Context, userId int64) ([]dto.File, error)
	SetFileStatus(ctx context.Context, id int64, status dto.FileStatus) error
}
