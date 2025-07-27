package main

import (
	"context"
	qrv1 "github.com/mdqni/Attendly/proto/gen/go/qr/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authv1 "github.com/mdqni/Attendly/proto/gen/go/auth/v1"
	groupv1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user:50051", opts); err != nil {
		log.Fatalf("failed to register user-service: %v", err)
	}

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "auth:50050", opts); err != nil {
		log.Fatalf("failed to register auth-service: %v", err)
	}

	if err := groupv1.RegisterGroupServiceHandlerFromEndpoint(ctx, mux, "group:50052", opts); err != nil {
		log.Fatalf("failed to register group-service: %v", err)
	}

	if err := qrv1.RegisterQRServiceHandlerFromEndpoint(ctx, mux, "qr:50054", opts); err != nil {
		log.Fatalf("failed to register group-service: %v", err)
	}

	log.Println("🚀 HTTP API Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
