package grpc

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type userServer struct {
	userv1.UnimplementedUserServiceServer
	service service.UserService
}
type UserService interface {
	Register(ctx context.Context, name, barcode, role string) (*userv1.User, error)
	GetUser(ctx context.Context, id string) (*userv1.User, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
	Login(ctx context.Context, name, password string) (*userv1.User, error)
	HasPermission(ctx context.Context, userID, permission string) (bool, error)
}

func (h *userServer) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	has, err := h.service.HasPermission(ctx, userID, permission)
	if err != nil {
		return false, err
	}
	return has, nil
}

func Register(gRPCServer *grpc.Server, svc service.UserService) {
	userv1.RegisterUserServiceServer(gRPCServer, &userServer{service: svc})
}

func (h *userServer) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := h.service.Login(ctx, request.Barcode, request.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &userv1.LoginResponse{Token: user.GetToken(), User: user.GetUser()}, nil
}

func (h *userServer) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	user, err := h.service.Register(ctx, req.GetName(), req.GetBarcode(), req.GetPassword(), req.GetRole())

	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{User: user}, nil
}

func (h *userServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if strings.TrimSpace(req.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "no user ID")
	}
	user, err := h.service.GetUserById(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, "User not found")
	}
	return &userv1.GetUserResponse{User: sanitizeUser(user)}, nil
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

func sanitizeUser(u *userv1.User) *userv1.User {
	u.Password = ""
	return u
}
