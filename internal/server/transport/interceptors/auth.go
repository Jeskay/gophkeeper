package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	proto "gophkeeper/api/protos"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/dto"
)

func NewAuthUnaryInterceptor(auth auth.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if _, ok := req.(*proto.RegisterRequest); ok {
			resp, err = handler(ctx, req)
			return
		}
		if _, ok := req.(*proto.LoginRequest); ok {
			resp, err = handler(ctx, req)
			return
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		jwtHeader, ok := md["jwt"]
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing jwt token")
		}
		token := jwtHeader[0]
		user, err := auth.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		resp, err = handler(context.WithValue(ctx, dto.Id, user.Id), req)
		return
	}
}

func NewAuthStreamInterceptor(auth auth.Service) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok || len(md["jwt"]) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing auth token")
		}
		token := md["jwt"][0]
		user, err := auth.VerifyToken(token)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		newCtx := context.WithValue(ss.Context(), dto.Id, user.Id)
		return handler(srv, &wrappedAuthStream{ServerStream: ss, WrappedContext: newCtx})
	}
}

type wrappedAuthStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

func (w *wrappedAuthStream) Context() context.Context {
	return w.WrappedContext
}
