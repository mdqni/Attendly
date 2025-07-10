package repository

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*userv1.User, error)
	GetAllUsers(ctx context.Context, page int, limit int) ([]*userv1.User, error)
	UpdateUser(ctx context.Context, u *userv1.User) (*userv1.User, error)
	DeleteUser(ctx context.Context, id string) error
	IsUserInGroup(ctx context.Context, userID string, groupID string) (bool, error)
}
