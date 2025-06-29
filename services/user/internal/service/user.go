package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/config"
	"github.com/mdqni/Attendly/services/user/internal/repository"
	passwordUtils "github.com/mdqni/Attendly/services/user/internal/utils/passwordUtils"
	"github.com/mdqni/Attendly/services/user/internal/utils/token"
	errs "github.com/mdqni/Attendly/shared/errs"
	"github.com/mdqni/Attendly/shared/rate_limit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type userService struct {
	repo    repository.UserRepository
	limiter rate_limit.LimiterInterface
	cfg     *config.Config
}

func NewUserService(repo repository.UserRepository, limiter rate_limit.LimiterInterface, cfg *config.Config) UserService {
	if limiter == nil {
		panic("rate limiter is nil")
	}
	return &userService{repo: repo, limiter: limiter, cfg: cfg}
}

func (s *userService) Register(ctx context.Context, name, barcode, password, role string) (*userv1.User, error) {
	if name == "" || barcode == "" || password == "" || role == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrMissingField.Error())
	}
	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, errs.ErrPasswordTooShort.Error())
	}

	exists, err := s.repo.GetUserByBarcode(ctx, barcode)
	if err != nil {
		log.Println("error checking user existence:", err)
		return nil, status.Error(codes.Internal, "failed to check user")
	}
	if exists != nil {
		return nil, status.Error(codes.AlreadyExists, errs.ErrUserAlreadyExists.Error())
	}

	hashed, err := passwordUtils.HashPassword(password)
	if err != nil {
		log.Println("failed to hash password:", err)
		return nil, status.Error(codes.Internal, "failed to process password")
	}

	user := &userv1.User{
		Id:       uuid.NewString(),
		Name:     name,
		Barcode:  barcode,
		Password: hashed,
		Role:     role,
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		log.Println("failed to save user:", err)
		return nil, status.Error(codes.Internal, errs.ErrOnUserSaving.Error())
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, barcode string, password string) (*userv1.LoginResponse, error) {
	if s.limiter == nil {
		return nil, status.Error(codes.Internal, "rate limiter not configured")
	}

	key := "login:" + barcode
	allowed, err := s.limiter.Allow(ctx, key, 5, time.Minute)
	if err != nil {
		log.Println("rate limiter error:", err)
		return nil, status.Error(codes.Internal, "rate limiter error")
	}
	if !allowed {
		return nil, status.Error(codes.ResourceExhausted, errs.ErrTooManyLoginAttempt.Error())
	}

	user, err := s.repo.GetUserByBarcode(ctx, barcode)
	if errors.Is(err, errs.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	log.Println(user.Password)
	log.Println(password)
	if err := passwordUtils.CheckPassword(user.Password, password); err != nil {
		log.Println("invalid password:", err)
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	perms, err := s.repo.GetPermissions(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not fetch permissions")
	}

	tokenStr, err := token.GenerateJWT(s.cfg.JwtSecret, user.Id, perms, time.Hour*24)
	if err != nil {
		return nil, status.Error(codes.Internal, "token generation failed")
	}
	if err := s.limiter.Reset(ctx, key); err != nil {
		log.Println("limiter reset error:", err)
	}

	return &userv1.LoginResponse{
		Token: tokenStr,
		User:  user,
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *userService) GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error) {
	return s.repo.GetUserByBarcode(ctx, barcode)
}

func (s *userService) IsInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return s.repo.CheckUserInGroup(ctx, userID, groupID)
}

func (s *userService) HasPermission(ctx context.Context, userID, action string) (bool, error) {
	return s.repo.HasPermission(ctx, userID, action)
}
