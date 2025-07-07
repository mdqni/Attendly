package service

import (
	"context"

	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*userv1.User, error)
	GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error)
	GetAllUsers(ctx context.Context) ([]*userv1.User, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
}
