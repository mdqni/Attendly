package grpc

import (
	"context"
	"log"
	"strings"

	group "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	"github.com/mdqni/Attendly/services/group/internal/service"
	errs "github.com/mdqni/Attendly/shared/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type groupServer struct {
	group.UnimplementedGroupServiceServer
	service service.GroupService
}

func (g *groupServer) CreateGroup(ctx context.Context, req *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	if strings.TrimSpace(req.GroupName) == "" || strings.TrimSpace(req.Department) == "" || req.Year == 0 {
		return nil, status.Error(codes.InvalidArgument, errs.ErrMissingField.Error())
	}

	grp, err := g.service.CreateGroup(ctx, req.GroupName, req.Department, int(req.Year))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &group.CreateGroupResponse{GroupId: grp.Id}, nil
}

func (g *groupServer) AddUserToGroup(ctx context.Context, req *group.AddUserToGroupRequest) (*group.AddUserToGroupResponse, error) {
	if strings.TrimSpace(req.GroupId) == "" || strings.TrimSpace(req.UserId) == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrMissingField.Error())
	}

	success, err := g.service.AddUserToGroup(ctx, req.GroupId, req.UserId)
	if err != nil {
		log.Println("AddUserToGroup error:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &group.AddUserToGroupResponse{Success: success}, nil
}

func (g *groupServer) RemoveUserFromGroup(ctx context.Context, req *group.RemoveUserFromGroupRequest) (*group.RemoveUserFromGroupResponse, error) {
	success, err := g.service.RemoveUserFromGroup(ctx, req.GroupId, req.UserId)
	if err != nil {
		log.Println("RemoveUserFromGroup error:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &group.RemoveUserFromGroupResponse{Success: success}, nil
}

func (g *groupServer) GetGroup(ctx context.Context, req *group.GetGroupRequest) (*group.GetGroupResponse, error) {
	grp, err := g.service.GetGroup(ctx, req.GroupId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &group.GetGroupResponse{Group: grp}, nil
}

func (g *groupServer) ListUsersInGroup(ctx context.Context, req *group.ListUsersInGroupRequest) (*group.ListUsersInGroupResponse, error) {
	users, err := g.service.ListUsersInGroup(ctx, req.GroupId)
	if err != nil {
		log.Println("ListUsersInGroup error:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &group.ListUsersInGroupResponse{User: users}, nil
}

func Register(gRPCServer *grpc.Server, svc service.GroupService) {
	if gRPCServer == nil || svc == nil {
		log.Fatal("Register: gRPC server or service is nil")
	}
	group.RegisterGroupServiceServer(gRPCServer, &groupServer{service: svc})
}
