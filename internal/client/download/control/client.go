package control

import (
	"context"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/client/download/abstraction"
	"strconv"
)

type downloadClient struct {
	grpcClient proto.GophkeeperClient
}

type downloadStream struct {
	stream proto.Gophkeeper_DownloadFileClient
}

func NewClient(grpcClient proto.GophkeeperClient) *downloadClient {
	return &downloadClient{grpcClient: grpcClient}
}

func NewStream(ctx context.Context, grpcClient proto.GophkeeperClient, req *proto.DownloadRequest) (*downloadStream, error) {
	ss, err := grpcClient.DownloadFile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &downloadStream{stream: ss}, nil
}

func (c *downloadClient) GetFiles(ctx context.Context) ([]*abstraction.FileData, error) {
	res, err := c.grpcClient.GetFiles(ctx, &proto.GetRequest{})
	if err != nil {
		return nil, err
	}
	fd := make([]*abstraction.FileData, len(res.Files))
	for i, f := range res.Files {
		size := strconv.FormatUint(uint64(f.GetSize()), 10)
		status := f.GetStatus().String()
		fd[i] = &abstraction.FileData{Name: f.Name + f.FileType, Id: f.GetId(), Size: size, Status: status}
	}
	return fd, nil
}

func (c *downloadClient) StartDownload(ctx context.Context, id int64) (abstraction.DownloadStream, error) {
	req := &proto.DownloadRequest{
		Id: id,
	}
	return NewStream(ctx, c.grpcClient, req)
}

func (s *downloadStream) Init() (string, error) {
	res, err := s.stream.Recv()
	if err != nil {
		return "", err
	}
	info := res.GetInfo()
	return info.Name + info.FileType, nil
}

func (s *downloadStream) Receive() ([]byte, error) {
	res, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	return res.GetChunkData(), nil
}
