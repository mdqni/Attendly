package service

import (
	"context"
	groupv1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type GroupService interface {
	CreateGroup(ctx context.Context, groupName string, department string, year int) (*groupv1.Group, error)
	AddUserToGroup(ctx context.Context, groupId string, userId string) (bool, error)
	RemoveUserFromGroup(ctx context.Context, groupId string, userId string) (bool, error)
	GetGroup(ctx context.Context, groupID string) (*groupv1.Group, error)
	ListUsersInGroup(ctx context.Context, groupId string) ([]*userv1.User, error)
	IsInGroup(ctx context.Context, groupId string, userId string) (bool, error)
}
