package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mdqni/Attendly/services/auth/internal/config"
	"github.com/mdqni/Attendly/services/auth/internal/domain/model"
	"github.com/mdqni/Attendly/services/auth/internal/repository"
	"github.com/mdqni/Attendly/shared/errs"
	"github.com/mdqni/Attendly/shared/passwordUtils"
	"github.com/mdqni/Attendly/shared/redisUtils"
	"github.com/mdqni/Attendly/shared/token"
	"time"
)

type authService struct {
	repo    repository.AuthRepository
	limiter redisUtils.LimiterInterface
	cfg     *config.Config
}

func NewAuthService(repo repository.AuthRepository, limiter redisUtils.LimiterInterface, cfg *config.Config) AuthService {
	return &authService{repo: repo, limiter: limiter, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	if input.Name == "" || input.Barcode == "" || input.Password == "" || input.Role == "" {
		return nil, errs.ErrMissingField
	}
	if len(input.Password) < 8 {
		return nil, errs.ErrPasswordTooShort
	}

	exists, err := s.repo.GetUserByBarcode(ctx, input.Barcode)
	if err == nil && exists != nil {
		return nil, errs.ErrUserAlreadyExists
	}

	hashed, err := passwordUtils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := model.UserWithPassword{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Barcode:  input.Barcode,
		Password: hashed,
		Role:     input.Role,
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	perms, err := s.repo.GetPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.GenerateJWT(s.cfg.JwtSecret, user.ID, perms, time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	if input.Barcode == "" || input.Password == "" {
		return nil, errs.ErrMissingField
	}

	key := "login:" + input.Barcode
	allowed, err := s.limiter.Allow(ctx, key, 5, time.Minute)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errs.ErrTooManyLoginAttempt
	}

	user, err := s.repo.GetUserByBarcode(ctx, input.Barcode)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	if err := passwordUtils.CheckPassword(user.Password, input.Password); err != nil {
		return nil, errs.ErrInvalidPassword
	}

	perms, err := s.repo.GetPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.GenerateJWT(s.cfg.JwtSecret, user.ID, perms, time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
