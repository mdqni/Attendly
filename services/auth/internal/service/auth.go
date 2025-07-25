package service

import (
	"context"
	"github.com/google/uuid"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/auth/internal/config"
	"github.com/mdqni/Attendly/services/auth/internal/domain/model"
	kafka2 "github.com/mdqni/Attendly/services/auth/internal/kafka"
	"github.com/mdqni/Attendly/services/auth/internal/repository"
	"github.com/mdqni/Attendly/shared/errs"
	"github.com/mdqni/Attendly/shared/passwordUtils"
	"github.com/mdqni/Attendly/shared/redisUtils"
	"github.com/mdqni/Attendly/shared/token"
	"log"
	"time"
)

type authService struct {
	repo          repository.AuthRepository
	limiter       redisUtils.LimiterInterface
	kafkaProducer *kafka2.EventProducer
	cfg           *config.Config
}

func NewAuthService(repo repository.AuthRepository, limiter redisUtils.LimiterInterface, cfg *config.Config, kafka *kafka2.EventProducer) AuthService {
	return &authService{repo: repo, limiter: limiter, cfg: cfg, kafkaProducer: kafka}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	const op = "service.auth.register"
	log.Println("op", op)

	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	user := createUserFromInput(input)
	if err := s.repo.SaveUser(ctx, user); err != nil {
		log.Println("op", op, "SaveUser", err)
		return nil, err
	}

	perms, err := s.repo.GetPermissions(ctx, user.ID)
	if err != nil {
		log.Println("op", op, "GetPermissions", err)
		return nil, err
	}

	accessToken, err := token.GenerateJWT(s.cfg.JwtSecret, user.ID, input.Role, perms, time.Hour)
	if err != nil {
		log.Println("op", op, "GenerateJWT", err)
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		log.Println("op", op, "GenerateRefreshToken", err)
		return nil, err
	}

	if err := s.kafkaProducer.SendUserRegisteredEvent(ctx, user.ID, user.Email, user.Role, user.Name); err != nil {
		log.Println("op", op, "Kafka send failed", err)
	}

	log.Println("Register success:", user.ID)
	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: userv1.User{
			Id:      user.ID,
			Name:    user.Name,
			Barcode: user.Barcode,
			Role:    user.Role,
			Email:   user.Email,
		},
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

	accessToken, err := token.GenerateJWT(s.cfg.JwtSecret, user.ID, user.Role, perms, time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User: userv1.User{
			Id:        user.ID,
			Name:      user.Name,
			Barcode:   user.Barcode,
			Role:      user.Role,
			Email:     user.Email,
			AvatarUrl: nil,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validateRegisterInput(input RegisterInput) error {
	if input.Name == "" || input.Barcode == "" || input.Password == "" || input.Role == "" {
		return errs.ErrMissingField
	}
	if len(input.Password) < 8 {
		return errs.ErrPasswordTooShort
	}
	return nil
}

func createUserFromInput(input RegisterInput) model.UserWithPassword {
	hashed, _ := passwordUtils.HashPassword(input.Password)

	return model.UserWithPassword{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Barcode:  input.Barcode,
		Email:    input.Email,
		Password: hashed,
		Role:     input.Role,
	}
}
