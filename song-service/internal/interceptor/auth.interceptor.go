package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("x-user-id")
			if len(values) > 0 {
				ctx = context.WithValue(ctx, UserIDKey, strings.TrimSpace(values[0]))
			}
		}

		return handler(ctx, req)
	}
}
