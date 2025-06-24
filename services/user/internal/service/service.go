package service

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type UserService interface {
	Register(ctx context.Context, name, barcode, password, role string) (*userv1.User, error)
	GetUserById(ctx context.Context, id string) (*userv1.User, error)
	GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error)
	Login(ctx context.Context, name, password string) (*userv1.LoginResponse, error)
	IsInGroup(ctx context.Context, userID, groupID string) (bool, error)
	HasPermission(ctx context.Context, userID, action string) (bool, error)
}
