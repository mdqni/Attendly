package interceptor

import (
	"context"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

func RBACInterceptor(svc service.UserService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID, ok := UserIDFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no user ID")
		}

		action := normalizeAction(info.FullMethod)
		log.Println("RBAC interceptor:", userID, action)

		ok, err := svc.HasPermission(ctx, userID, action)
		if err != nil {
			return nil, status.Error(codes.Internal, "permission check error")
		}
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "no permission")
		}

		return handler(ctx, req)
	}
}

func normalizeAction(fullMethod string) string {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	fullMethod = strings.ReplaceAll(fullMethod, ".", "_")
	fullMethod = strings.ReplaceAll(fullMethod, "/", "_")
	return toSnakeCase(fullMethod)
}

func toSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
