package transport

import (
	"context"
	auth "gophkeeper/api/protos"
	"gophkeeper/internal/server/endpoints"
)

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.LoginRequest)
	return endpoints.LoginRequest{Name: req.Name, Password: req.Password}, nil
}

func decodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.RegisterRequest)
	return endpoints.RegisterRequest{Name: req.Name, Password: req.Password}, nil
}
func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.LoginResponse)
	return &auth.LoginResponse{Status: response.Status, Token: response.Token}, nil
}

func encodeGRPCRegisterResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	_ = grpcResponse.(endpoints.RegisterResponse)
	return &auth.RegisterResponse{}, nil
}
