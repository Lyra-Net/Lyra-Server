package interceptor

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserAgentKey contextKey = "user_agent"
	UserIpKey    contextKey = "user_ip"
)

var metaToContextKeys = map[string]contextKey{
	"x-user-id":       UserIDKey,
	"user-agent":      UserAgentKey,
	"x-forwarded-for": UserIpKey,
}

func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		log.Println("metadata: ", md)
		if ok {
			for metaKey, ctxKey := range metaToContextKeys {
				values := md.Get(metaKey)
				if len(values) > 0 {
					ctx = context.WithValue(ctx, ctxKey, strings.TrimSpace(values[0]))
				}
			}
		}

		return handler(ctx, req)
	}
}
