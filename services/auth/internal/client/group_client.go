package client

import (
	"context"
	"fmt"
	group1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	"github.com/mdqni/Attendly/shared/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GroupClient struct {
	conn   *grpc.ClientConn
	client group1.GroupServiceClient
}

func NewGroupClient(addr string) (*GroupClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor.UnaryAuthForwardInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to group service: %w", err)
	}
	return &GroupClient{
		conn:   conn,
		client: group1.NewGroupServiceClient(conn),
	}, nil
}

// Пока не нужен
func (g *GroupClient) IsInGroup(ctx context.Context, userId string, groupId string) (bool, error) {
	isInGroup, err := g.client.IsInGroup(ctx, &group1.IsInGroupRequest{
		GroupId: groupId,
		UserId:  userId,
	})
	if err != nil {
		return false, err
	}
	return isInGroup.IsMember, nil
}

func (g *GroupClient) Close() error {
	return g.conn.Close()
}
