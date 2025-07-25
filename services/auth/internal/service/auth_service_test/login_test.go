package service_test

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/auth/internal/config"
	"github.com/mdqni/Attendly/services/auth/internal/repository/mocks"
	"github.com/mdqni/Attendly/services/auth/internal/service"
	errs "github.com/mdqni/Attendly/shared/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

type fakeLimiter struct {
	l *redis_rate.Limiter
}

func (f *fakeLimiter) Allow(ctx context.Context, key string, rate int, period time.Duration) (bool, error) {
	return true, nil
}
func (f *fakeLimiter) Reset(_ context.Context, _ string) error {
	return nil
}
func TestLogin_Success(t *testing.T) {
	mockRepo := mocks.NewAuthRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewAuthService(mockRepo, limiter, cfg)

	user := &userv1.User{
		Id:      "u-1",
		Name:    "Test",
		Barcode: "123456",
		Role:    "student",
	}
	var perms []string

	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(user, nil)
	mockRepo.On("GetPermissions", mock.Anything, "u-1").Return(perms, nil)

	resp, err := svc.Login(context.Background(), service.LoginInput{
		Barcode:  "123456",
		Password: "1234",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.Id, resp.UserID)
	assert.NotEmpty(t, resp.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {

	mockRepo := mocks.NewAuthRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewAuthService(mockRepo, limiter, cfg)
	user := &userv1.User{
		Id:      "u-1",
		Name:    "Test",
		Barcode: "123456",
		Role:    "student",
	}
	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(user, nil)

	resp, err := svc.Login(context.Background(), service.LoginInput{
		Barcode: "123456", Password: "wrong-password",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := mocks.NewAuthRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewAuthService(mockRepo, limiter, cfg)

	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(nil, errs.ErrUserNotFound)
	resp, err := svc.Login(context.Background(), service.LoginInput{Barcode: "123456", Password: "1234"})
	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	mockRepo.AssertExpectations(t)
}

type limitedFakeLimiter struct {
	attempts int
	limit    int
}

func (l *limitedFakeLimiter) Allow(ctx context.Context, key string, rate int, period time.Duration) (bool, error) {
	l.attempts++
	return l.attempts <= l.limit, nil
}

func (l *limitedFakeLimiter) Reset(_ context.Context, _ string) error {
	l.attempts = 0
	return nil
}

func TestLogin_ResourceExhausted(t *testing.T) {
	mockRepo := mocks.NewAuthRepository(t)
	limiter := &limitedFakeLimiter{limit: 5}
	cfg := config.MustLoad()

	svc := service.NewAuthService(mockRepo, limiter, cfg)
	user := &userv1.User{
		Id:      "u-1",
		Name:    "Test",
		Barcode: "123456",
		Role:    "student",
	}
	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(user, nil)
	resp, err := svc.Login(context.Background(), service.LoginInput{Barcode: "123456", Password: "12347456"})

	for i := 0; i < 5; i++ {
		resp, err = svc.Login(context.Background(), service.LoginInput{Barcode: "123456", Password: "12345644"})
	}
	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.ResourceExhausted, st.Code())
	mockRepo.AssertExpectations(t)
}
