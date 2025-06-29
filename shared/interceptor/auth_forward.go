package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryAuthForwardInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctxWithAuth := metadata.AppendToOutgoingContext(ctx, "authorization", tokens[0])
		return invoker(ctxWithAuth, method, req, reply, cc, opts...)
	}
}
