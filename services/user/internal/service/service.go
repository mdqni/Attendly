package service

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserService interface {
	Register(ctx context.Context, name, barcode, role string) (*userv1.User, error)
	GetUser(ctx context.Context, id string) (*userv1.User, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
	HasPermission(ctx context.Context, userID, action string) (bool, error)
}
