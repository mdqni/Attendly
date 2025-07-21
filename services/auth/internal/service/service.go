package service

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
}

type RegisterInput struct {
	Name     string
	Barcode  string
	Password string
	Email    string
	Role     string
}

type LoginInput struct {
	Barcode  string
	Password string
}

type AuthResult struct {
	User         userv1.User
	AccessToken  string
	RefreshToken string
}
