package client

import (
	"fmt"
	authv1 "github.com/mdqni/Attendly/proto/gen/go/auth/v1"
	"github.com/mdqni/Attendly/shared/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor.UnaryAuthForwardInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	return &AuthClient{
		conn:   conn,
		client: authv1.NewAuthServiceClient(conn),
	}, nil
}

func (g *AuthClient) Close() error {
	return g.conn.Close()
}
