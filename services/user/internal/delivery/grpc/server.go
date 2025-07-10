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

func Register(gRPCServer *grpc.Server, svc service.UserService) {
	userv1.RegisterUserServiceServer(gRPCServer, &userServer{service: svc})
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
	return &userv1.GetUserResponse{User: user}, nil
}
func (h *userServer) GetAllUsers(ctx context.Context, req *userv1.GetUsersRequest) (userv1.GetUsersResponse, error) {
	user, err := h.service.GetAllUsers(ctx, int(req.GetPage()), int(req.GetLimit()))
	if err != nil {
		log.Println(err)
		return userv1.GetUsersResponse{
			Users: nil,
		}, status.Error(codes.NotFound, "Users not found")
	}
	return userv1.GetUsersResponse{Users: user}, nil
}

func (h *userServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if strings.TrimSpace(req.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "no user ID")
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}
