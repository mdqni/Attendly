package service

import (
	"context"
	"errors"
	groupv1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/group/internal/client"
	"github.com/mdqni/Attendly/services/group/internal/config"
	"github.com/mdqni/Attendly/services/group/internal/repository"
	"github.com/mdqni/Attendly/shared/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type groupService struct {
	repo       repository.GroupRepository
	cfg        *config.Config
	userClient *client.UserClient
}

func (s *groupService) IsInGroup(ctx context.Context, groupId string, userId string) (bool, error) {
	inGroup, err := s.repo.IsInGroup(ctx, groupId, userId)
	if err != nil {
		return false, err
	}
	return inGroup, nil
}

func (s *groupService) CreateGroup(ctx context.Context, groupName string, department string, year int) (*groupv1.Group, error) {
	group, err := s.repo.CreateGroup(ctx, groupName, department, year)

	if err != nil {
		if errors.Is(err, errs.ErrGroupAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "group already exists")
		}
		return nil, err
	}
	return group, nil
}

func (s *groupService) GetGroup(ctx context.Context, groupID string) (*groupv1.Group, error) {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		if errors.Is(err, errs.ErrGroupNotFound) {
			return nil, status.Error(codes.NotFound, "group not found")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return group, nil
}

func (s *groupService) AddUserToGroup(ctx context.Context, groupId string, userId string) (bool, error) {
	_, err := s.userClient.GetUser(ctx, userId)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return false, status.Error(codes.NotFound, "user not found")
		}
		log.Println(err)
		return false, err

	}
	group, err := s.repo.AddUserToGroup(ctx, groupId, userId)
	if err != nil {
		if errors.Is(err, errs.ErrUserOrGroupNotExists) {
			return false, status.Error(codes.NotFound, "user or group not exists")
		}
		return false, status.Error(codes.Internal, "error adding user to group")
	}
	return group, nil
}

func (s *groupService) RemoveUserFromGroup(ctx context.Context, groupId string, userId string) (bool, error) {
	isRemoved, err := s.repo.RemoveUserFromGroup(ctx, groupId, userId)
	if err != nil {
		if errors.Is(err, errs.ErrUserOrGroupNotExists) {
			return false, status.Error(codes.NotFound, "user or group not exists")
		}
		log.Println(err)
		return false, status.Error(codes.Internal, "Internal error")
	}
	return isRemoved, nil
}

func (s *groupService) ListUsersInGroup(ctx context.Context, groupID string) ([]*userv1.User, error) {
	users, err := s.repo.ListUsersInGroup(ctx, groupID)
	if err != nil {
		if errors.Is(err, errs.ErrGroupIsEmpty) {
			return nil, status.Error(codes.NotFound, "group is empty")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return users, nil
}

func NewGroupService(repo repository.GroupRepository, cfg *config.Config, userClient *client.UserClient) GroupService {
	return &groupService{repo: repo, cfg: cfg, userClient: userClient}
}
