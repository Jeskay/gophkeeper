package abstraction

import (
	"context"
	proto "gophkeeper/api/protos"
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

func (c *downloadClient) GetFiles(ctx context.Context) ([]*FileData, error) {
	res, err := c.grpcClient.GetFiles(ctx, &proto.GetRequest{})
	if err != nil {
		return nil, err
	}
	fd := make([]*FileData, len(res.Files))
	for i, f := range res.Files {
		fd[i] = &FileData{Name: f} //TODO: provide all file information from server
	}
	return fd, nil
}

func (c *downloadClient) StartDownload(ctx context.Context, fileName string) (DownloadStream, error) {
	req := &proto.DownloadRequest{
		Name: fileName,
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

func (s *downloadStream) Close() error {
	return s.stream.CloseSend()
}
