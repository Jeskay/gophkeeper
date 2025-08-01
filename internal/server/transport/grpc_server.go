package transport

import (
	"bytes"
	"context"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/db"
	"gophkeeper/internal/server/dto"
	"gophkeeper/internal/server/file"
	"io"
)

type KeeperServer struct {
	authService auth.Service
	fileService file.Service
	dbService   db.Service
	proto.UnimplementedGophkeeperServer
}

func NewGRPCServer(authService auth.Service, fileService file.Service, dbService db.Service) proto.GophkeeperServer {
	return &KeeperServer{
		authService: authService,
		fileService: fileService,
		dbService:   dbService,
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

	files, err := k.dbService.GetUserFiles(ctx, userId)
	if err != nil {
		return &proto.GetResponse{}, err
	}
	fNames := make([]string, len(files))
	for i, f := range files {
		fNames[i] = f.Name
	}
	return &proto.GetResponse{Files: fNames}, nil
}

func (k *KeeperServer) SaveFile(stream proto.Gophkeeper_SaveFileServer) error {
	v := stream.Context().Value(dto.Id)
	userId := v.(int64)

	req, err := stream.Recv()
	if err != nil {
		return err
	}
	fName := req.GetInfo().GetName()
	fType := req.GetInfo().GetFileType()

	fId, err := k.dbService.CreateFile(context.Background(), dto.File{Name: fName + fType, Status: dto.Processing}, userId)
	if err != nil {
		return err
	}
	data := bytes.Buffer{}
	var size int
	for {
		if stream.Context().Err() == context.Canceled {
			return context.Canceled
		}
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunk := req.GetChunkData()
		size, err = data.Write(chunk)
		if err != nil {
			return err
		}
	}
	err = stream.SendAndClose(&proto.SaveResponse{Size: uint32(size)})
	if err != nil {
		return err
	}
	if err := k.fileService.SaveFile(fName+fType, data.Bytes()); err != nil {
		return err
	}
	err = k.dbService.SetFileStatus(stream.Context(), fId, dto.Available)
	if err != nil {
		return err
	}
	return nil
}

func (k *KeeperServer) DownloadFile(req *proto.DownloadRequest, stream proto.Gophkeeper_DownloadFileServer) error {
	v := stream.Context().Value(dto.Id)
	userId := v.(int64)

	fInfo, err := k.dbService.GetFile(stream.Context(), userId, req.Name)
	if err != nil {
		return err
	}
	err = k.dbService.SetFileStatus(stream.Context(), fInfo.Id, dto.Reserved)
	if err != nil {
		return nil
	}
	defer k.dbService.SetFileStatus(stream.Context(), fInfo.Id, dto.Available)
	stream.Send(&proto.DownloadResponse{Data: &proto.DownloadResponse_Info{Info: &proto.FileInfo{Name: fInfo.Name}}})
	return k.fileService.ReadByChunk(fInfo.Name, func(b []byte) error {
		return stream.Send(&proto.DownloadResponse{Data: &proto.DownloadResponse_ChunkData{ChunkData: b}})
	})
}
