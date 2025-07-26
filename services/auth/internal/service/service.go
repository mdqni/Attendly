package service

import (
	"context"
	authv1 "github.com/mdqni/Attendly/proto/gen/go/auth/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	Refresh(ctx context.Context, req *authv1.RefreshTokenRequest) (*AuthResult, error)
	GetUserInfoById(ctx context.Context, id string) (*userv1.User, error)
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
