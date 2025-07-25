package interceptors

import (
	"context"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/dto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthUnaryInterceptor(auth auth.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		token := md["jwt"][0]
		user, err := auth.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		resp, err = handler(context.WithValue(ctx, dto.Id, user.Id), req)
		return
	}
}
