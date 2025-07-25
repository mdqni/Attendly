package interceptor

import (
	"context"
	"github.com/mdqni/Attendly/shared/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

var openMethods = map[string]struct{}{
	"/auth.v1.AuthService/Login":    {},
	"/auth.v1.AuthService/Register": {},
}

const userIDKey = "user_id_key"

func RBACInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println(info.FullMethod)

		if _, ok := openMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")
		claims, err := token.ParseJWT(tokenStr, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		if claims.Role == "admin" {
			return handler(ctx, req)
		}
		action := normalizeAction(info.FullMethod)
		log.Println(action)
		if !contains(claims.Perms, action) {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		ctx = context.WithValue(ctx, userIDKey, claims.UserID)
		return handler(ctx, req)
	}
}
func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func normalizeAction(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	fullMethod = strings.ReplaceAll(fullMethod, ".", "_")
	fullMethod = strings.ReplaceAll(fullMethod, "/", "_")
	return fullMethod
}
