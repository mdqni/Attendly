package repository

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *userv1.User) error
	GetUserById(ctx context.Context, barcode string) (*userv1.User, error)
	GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error)
	HasPermission(ctx context.Context, userID, action string) (bool, error)
	GetPermissions(ctx context.Context, userID string) ([]string, error)
}
