package keeper

import (
	"context"
	"fmt"

	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/server/file"
)

type keeperService struct {
	dbService   db.Service
	fileService file.Service
}

func NewKeeperService(dbService db.Service, fileService file.Service) *keeperService {
	return &keeperService{dbService: dbService, fileService: fileService}
}

func (k *keeperService) GetFiles(ctx context.Context, userID int64) ([]dto.File, error) {
	files, err := k.dbService.GetUserFiles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (k *keeperService) UploadFile(ctx context.Context, userID int64, file dto.File, chData <-chan dto.ChunkData) (int, error) {
	file.Status = dto.Processing
	fID, err := k.dbService.CreateFile(ctx, file, userID)
	if err != nil {
		return 0, err
	}
	size, err := k.fileService.WriteByChunk(ctx, k.saveDir(userID, fID), chData)
	if err != nil {
		return 0, err
	}
	err = k.dbService.SetFileStatus(ctx, fID, dto.Available)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (k *keeperService) DownloadFile(ctx context.Context, userID int64, file dto.File, out chan<- dto.ChunkData) (*dto.File, error) {
	fInfo, err := k.dbService.GetFile(ctx, userID, file.Id)
	if err != nil {
		return &fInfo, err
	}
	if fInfo.Status != dto.Available {
		return &fInfo, fmt.Errorf("access restricted due to %s status of the file", fInfo.Status)
	}
	err = k.dbService.SetFileStatus(ctx, fInfo.Id, dto.Reserved)
	if err != nil {
		return &fInfo, err
	}
	go func() {
		defer close(out)
		if err := k.fileService.ReadByChunk(k.saveDir(userID, file.Id), out); err != nil {
			out <- dto.ChunkData{Err: err}
			return
		}
		if err := k.dbService.SetFileStatus(ctx, fInfo.Id, dto.Available); err != nil {
			out <- dto.ChunkData{Err: err}
		}
	}()
	return &fInfo, nil
}

func (k *keeperService) saveDir(userID int64, fileID int64) string {
	return fmt.Sprintf("%d-%d.dat", userID, fileID)
}
