package service

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
}

type RegisterInput struct {
	Name     string
	Barcode  string
	Password string
	Role     string
}

type LoginInput struct {
	Barcode  string
	Password string
}

type AuthResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
