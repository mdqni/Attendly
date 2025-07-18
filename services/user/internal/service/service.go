package service

import (
	"context"
	"github.com/mdqni/Attendly/shared/domain"

	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*userv1.User, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]*userv1.User, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
	CreateUser(ctx context.Context, user *domain.User) (*userv1.User, error)
}
