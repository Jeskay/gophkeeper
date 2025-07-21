package db

import (
	"context"
	"gophkeeper/internal/server/dto"
)

type DatabaseService interface {
	GetUser(ctx context.Context, name string) (dto.User, error)
	CreateUser(ctx context.Context, user dto.User) error
	GetUserFiles(ctx context.Context, userId int64) ([]dto.File, error)
	SetFileStatus(ctx context.Context, id int64, status dto.FileStatus) error
}
