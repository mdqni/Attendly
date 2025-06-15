package grpc

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	userv1.UnimplementedUserServiceServer
	service service.UserService
}
type UserService interface {
	Register(ctx context.Context, name, barcode, role string) (*userv1.User, error)
	GetUser(ctx context.Context, id string) (*userv1.User, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
}

func Register(gRPCServer *grpc.Server, service UserService) {
	userv1.RegisterUserServiceServer(gRPCServer, &userServer{service: service})
}

func (h *userServer) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	if req.GetRole() != "admin" || req.GetRole() != "teacher" {
		return nil, status.Error(codes.PermissionDenied, "role must be 'admin' or 'teacher'")
	}
	user, err := h.service.Register(ctx, req.GetName(), req.GetBarcode(), req.GetRole())

	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{User: user}, nil
}

func (h *userServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, req.GetId())
	if err != nil {

		return nil, status.Error(codes.NotFound, "User not found")
	}
	return &userv1.GetUserResponse{User: user}, nil
}

func (h *userServer) IsInGroup(ctx context.Context, req *userv1.IsInGroupRequest) (*userv1.IsInGroupResponse, error) {
	isInGroup, err := h.service.IsInGroup(ctx, req.GetUserId(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	return &userv1.IsInGroupResponse{
		IsMember: isInGroup,
	}, nil
}
