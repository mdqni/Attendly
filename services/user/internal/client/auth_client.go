package client

import (
	"context"
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

func (g *AuthClient) GetUserInfo(ctx context.Context, userId string) (*authv1.GetUserInfoResponse, error) {
	user, err := g.client.GetUserInfoById(ctx, &authv1.GetUserInfoRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return &authv1.GetUserInfoResponse{
		Name:    user.GetName(),
		Barcode: user.GetBarcode(),
		Role:    user.GetRole(),
	}, nil
}

func (g *AuthClient) Close() error {
	return g.conn.Close()
}
