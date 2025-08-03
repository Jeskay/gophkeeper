package transport

import (
	"context"
	"errors"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/server/keeper"
	"io"
)

type KeeperServer struct {
	authService   auth.Service
	keeperService keeper.Service
	proto.UnimplementedGophkeeperServer
}

func NewGRPCServer(authService auth.Service, keeperService keeper.Service) proto.GophkeeperServer {
	return &KeeperServer{
		authService:   authService,
		keeperService: keeperService,
	}
}

func (k *KeeperServer) Login(ctx context.Context, r *proto.LoginRequest) (*proto.LoginResponse, error) {
	status, token := k.authService.Login(ctx, r.Name, r.Password)
	return &proto.LoginResponse{Status: status, Token: token}, nil
}

func (k *KeeperServer) Register(ctx context.Context, r *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	err := k.authService.Register(ctx, r.Name, r.Password)
	return &proto.RegisterResponse{}, err
}

func (k *KeeperServer) GetFiles(ctx context.Context, r *proto.GetRequest) (*proto.GetResponse, error) {
	userId := ctx.Value(dto.Id).(int64)

	files, err := k.keeperService.GetFiles(ctx, userId)
	if err != nil {
		return &proto.GetResponse{}, err
	}
	fInfos := make([]*proto.FileInfo, len(files))
	for i, f := range files {
		fInfos[i] = f.Proto()
	}
	return &proto.GetResponse{Files: fInfos}, nil
}

func (k *KeeperServer) SaveFile(stream proto.Gophkeeper_SaveFileServer) (err error) {
	v := stream.Context().Value(dto.Id)
	userId := v.(int64)

	req, err := stream.Recv()
	if err != nil {
		return
	}
	var size int
	defer func() {
		closeErr := stream.SendAndClose(&proto.SaveResponse{Size: uint32(size)})
		if err == nil {
			err = closeErr
		}
	}()
	fInfo := req.GetInfo()
	if fInfo == nil {
		return errors.New("no file info provided")
	}
	fName := fInfo.GetName()
	fType := fInfo.GetFileType()
	fSize := fInfo.GetSize()

	file := dto.File{Name: fName + fType, Status: dto.Processing, Size: int(fSize)}
	chunkC := make(chan dto.ChunkData)
	go func() {
		defer close(chunkC)
		for {
			if stream.Context().Err() == context.Canceled {
				chunkC <- dto.ChunkData{Err: context.Canceled}
				return
			}
			req, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				chunkC <- dto.ChunkData{Err: err}
				break
			}
			chunkC <- dto.ChunkData{Data: req.GetChunkData()}
		}
	}()
	size, err = k.keeperService.UploadFile(context.Background(), userId, file, chunkC)
	return
}

func (k *KeeperServer) DownloadFile(req *proto.DownloadRequest, stream proto.Gophkeeper_DownloadFileServer) error {
	v := stream.Context().Value(dto.Id)
	userId := v.(int64)
	out := make(chan dto.ChunkData)
	file, err := k.keeperService.DownloadFile(stream.Context(), userId, dto.File{Id: req.Id}, out)
	if err != nil {
		return err
	}
	stream.Send(&proto.DownloadResponse{
		Data: &proto.DownloadResponse_Info{
			Info: &proto.FileInfo{
				Name: file.Name,
			},
		},
	})
	for chunk := range out {
		if chunk.Err != nil {
			return chunk.Err
		}
		err := stream.Send(&proto.DownloadResponse{
			Data: &proto.DownloadResponse_ChunkData{
				ChunkData: chunk.Data,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
