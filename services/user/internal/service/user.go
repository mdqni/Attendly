package service

import (
	"context"
	"github.com/google/uuid"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/repository"
	passwordUtils "github.com/mdqni/Attendly/services/user/internal/utils/passwordUtils"
	"github.com/mdqni/Attendly/services/user/internal/utils/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, name, barcode, password, role string) (*userv1.User, error) {
	user := &userv1.User{
		Id:       uuid.NewString(),
		Name:     name,
		Barcode:  barcode,
		Password: password,
		Role:     role,
	}

	err := s.repo.SaveUser(ctx, user)
	return user, err
}

func (s *userService) Login(ctx context.Context, barcode, password string) (*userv1.LoginResponse, error) {
	user, err := s.repo.GetUserByBarcode(ctx, barcode)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if err := passwordUtils.CheckPassword(user.Password, password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	perms, err := s.repo.GetPermissions(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not fetch permissions")
	}

	tokenStr, err := token.GenerateJWT("SUPER-SECRET-CODE", user.Id, perms, time.Hour*24)
	if err != nil {
		return nil, status.Error(codes.Internal, "token generation failed")
	}

	return &userv1.LoginResponse{
		Token: tokenStr,
		User:  user,
	}, nil
}
func (s *userService) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *userService) GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error) {
	return s.repo.GetUserById(ctx, barcode)
}

func (s *userService) IsInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return s.repo.CheckUserInGroup(ctx, userID, groupID)
}

func (s *userService) HasPermission(ctx context.Context, userID, action string) (bool, error) {
	permission, err := s.repo.HasPermission(ctx, userID, action)
	if err != nil {
		return false, err
	}
	return permission, nil
}
