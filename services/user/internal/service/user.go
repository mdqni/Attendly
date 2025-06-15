package service

import (
	"context"
	"github.com/google/uuid"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, name, barcode, role string) (*userv1.User, error) {
	user := &userv1.User{
		Id:      uuid.NewString(),
		Name:    name,
		Barcode: barcode,
		Role:    role,
	}

	err := s.repo.SaveUser(ctx, user)
	return user, err
}

func (s *userService) GetUser(ctx context.Context, id string) (*userv1.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) IsInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return s.repo.CheckUserInGroup(ctx, userID, groupID)
}
