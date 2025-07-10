package service

import (
	"context"
	"github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/client"
	"github.com/mdqni/Attendly/services/user/internal/config"
	"github.com/mdqni/Attendly/services/user/internal/repository"
)

type userService struct {
	repo        repository.UserRepository
	cfg         *config.Config
	groupClient *client.GroupClient
}

func NewUserService(repo repository.UserRepository, cfg *config.Config, group *client.GroupClient) UserService {
	return &userService{repo: repo, cfg: cfg, groupClient: group}
}

func (s *userService) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	byID, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &userv1.User{
		Id:      byID.Id,
		Name:    byID.Name,
		Barcode: byID.Barcode,
		Role:    byID.Role,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, page int, limit int) ([]*userv1.User, error) {
	return s.repo.GetAllUsers(ctx, page, limit)
}

func (s *userService) IsInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return s.repo.IsUserInGroup(ctx, userID, groupID)
}
