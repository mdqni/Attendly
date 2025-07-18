package service

import (
	"context"
	"errors"
	"github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/client"
	"github.com/mdqni/Attendly/services/user/internal/config"
	"github.com/mdqni/Attendly/services/user/internal/repository"
	"github.com/mdqni/Attendly/shared/domain"
	errPkg "github.com/mdqni/Attendly/shared/errs"
	"log"
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
		Id:        byID.Id,
		Name:      byID.Name,
		Email:     byID.Email,
		AvatarUrl: byID.AvatarUrl,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, page int, limit int) ([]*userv1.User, error) {
	return s.repo.GetAllUsers(ctx, page, limit)
}

func (s *userService) IsInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return s.repo.IsUserInGroup(ctx, userID, groupID)
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) (*userv1.User, error) {
	_, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, errPkg.ErrUserNotFound) {
			return s.repo.CreateUser(ctx, user)
		}
		log.Printf("failed to get user by id %s: %v", user.ID, err)
		return nil, err
	}

	log.Printf("user with id %s already exists", user.ID)
	return nil, errPkg.ErrUserAlreadyExists
}
