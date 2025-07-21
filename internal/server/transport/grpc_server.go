package transport

import (
	"context"
	proto "gophkeeper/api/protos"
	"gophkeeper/internal/server/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	login    grpckit.Handler
	register grpckit.Handler
	proto.UnimplementedGophkeeperServer
}

func NewGRPCServer(ep endpoints.Set) proto.GophkeeperServer {
	return &grpcServer{
		login:    grpckit.NewServer(ep.LoginEndpoint, decodeGRPCLoginRequest, encodeGRPCLoginResponse),
		register: grpckit.NewServer(ep.RegisterEndpoint, decodeGRPCRegisterRequest, encodeGRPCRegisterResponse),
	}
}

func (g *grpcServer) Login(ctx context.Context, r *proto.LoginRequest) (*proto.LoginResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.LoginResponse), nil
}

func (g *grpcServer) Register(ctx context.Context, r *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	_, resp, err := g.register.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.RegisterResponse), nil
}
