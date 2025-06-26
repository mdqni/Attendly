package service_test

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/repository/mocks"
	"github.com/mdqni/Attendly/services/user/internal/service"
	"github.com/mdqni/Attendly/services/user/internal/utils/passwordUtils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
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
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	svc := service.NewUserService(mockRepo, limiter)
	password, err := passwordUtils.HashPassword("1234")
	if err != nil {
		log.Println(err)
	}
	user := &userv1.User{
		Id:       "u-1",
		Name:     "Test",
		Barcode:  "123456",
		Password: password,
		Role:     "student",
	}
	perms := []string{"perm.view"}

	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(user, nil)
	mockRepo.On("GetPermissions", mock.Anything, "u-1").Return(perms, nil)

	resp, err := svc.Login(context.Background(), "123456", "1234")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.Id, resp.User.Id)
	assert.NotEmpty(t, resp.Token)
	mockRepo.AssertExpectations(t)
}
