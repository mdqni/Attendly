package repository

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *userv1.User) error
	GetUserByID(ctx context.Context, id string) (*userv1.User, error)
	CheckUserInGroup(ctx context.Context, userID, groupID string) (bool, error)
}
