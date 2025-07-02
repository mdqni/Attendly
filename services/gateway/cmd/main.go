package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := userv1.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"0.0.0.0:50051",
		opts,
	)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Println("🚀 HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
