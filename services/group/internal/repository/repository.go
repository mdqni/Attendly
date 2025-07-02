package repository

import (
	"context"
	groupv1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, groupName string, department string, year int) (*groupv1.Group, error)
	AddUserToGroup(ctx context.Context, groupID, userID string) (bool, error)
	RemoveUserFromGroup(ctx context.Context, groupID, userID string) (bool, error)
	GetGroup(ctx context.Context, groupID string) (*groupv1.Group, error)
	ListUsersInGroup(ctx context.Context, groupID string) ([]*userv1.User, error)
	IsInGroup(ctx context.Context, groupID string, userId string) (bool, error)
}
