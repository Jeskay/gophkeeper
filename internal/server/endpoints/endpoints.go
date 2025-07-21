package endpoints

import (
	"context"
	"gophkeeper/internal/server/auth"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	LoginEndpoint    endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	SaveFileEndpoint endpoint.Endpoint
}

func NewEndpointList(authSvc auth.Service) Set {
	return Set{
		LoginEndpoint:    MakeLoginEndpoint(authSvc),
		RegisterEndpoint: MakeRegisterEndpoint(authSvc),
	}
}

func MakeLoginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		status, token := svc.Login(ctx, req.Name, req.Password)
		return &LoginResponse{Status: status, Token: token}, nil
	}
}

func MakeRegisterEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		err := svc.Register(ctx, req.Name, req.Password)
		if err != nil {
			return &RegisterResponse{}, err
		}
		return &RegisterResponse{}, nil
	}
}
