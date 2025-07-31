package abstraction

import (
	"context"
	proto "gophkeeper/api/protos"
)

type uploadClient struct {
	grpcClient proto.GophkeeperClient
}

type uploadStream struct {
	stream proto.Gophkeeper_SaveFileClient
}

func NewClient(grpcClient proto.GophkeeperClient) *uploadClient {
	return &uploadClient{grpcClient: grpcClient}
}

func NewStream(ctx context.Context, grpcClient proto.GophkeeperClient) (*uploadStream, error) {
	stream, err := grpcClient.SaveFile(ctx)
	if err != nil {
		return nil, err
	}
	return &uploadStream{stream: stream}, nil
}

func (c *uploadClient) StartUpload(ctx context.Context) (UploadStream, error) {
	return NewStream(ctx, c.grpcClient)
}

func (s *uploadStream) Init(fileName, fileType string) error {
	req := &proto.SaveRequest{
		Data: &proto.SaveRequest_Info{
			Info: &proto.FileInfo{
				Name:     fileName,
				FileType: fileType,
			},
		},
	}
	return s.stream.Send(req)
}

func (s *uploadStream) Upload(data []byte) error {
	req := &proto.SaveRequest{
		Data: &proto.SaveRequest_ChunkData{
			ChunkData: data,
		},
	}
	return s.stream.Send(req)
}

func (s *uploadStream) Close() (uint32, error) {
	res, err := s.stream.CloseAndRecv()
	if err != nil || res == nil {
		return 0, err
	}
	return res.Size, nil
}
