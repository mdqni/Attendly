package repository

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/repository/postgres"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*postgres.InternalUser, error)
	GetUserByBarcode(ctx context.Context, barcode string) (*postgres.InternalUser, error)
	IsUserInGroup(ctx context.Context, userID string, groupID string) (bool, error)
	GetAllUsers(ctx context.Context) ([]*userv1.User, error)
}
