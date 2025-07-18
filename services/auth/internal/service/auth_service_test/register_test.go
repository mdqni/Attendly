package service_test

import (
	"context"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/config"
	"github.com/mdqni/Attendly/services/user/internal/repository/mocks"
	"github.com/mdqni/Attendly/services/user/internal/service"
	errs "github.com/mdqni/Attendly/shared/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)

	user := &userv1.User{
		Name:     "John Johnson",
		Barcode:  "123456",
		Password: "12345678",
		Role:     "student",
	}
	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(nil, nil)

	mockRepo.On("SaveUser", mock.Anything, mock.AnythingOfType("*userv1.User")).Return(nil)

	created, err := svc.Register(context.Background(), user.Name, user.Barcode, user.Password, user.Role)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, created.Name)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmptyFields(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)

	created, err := svc.Register(context.Background(), "", "", "", "")

	assert.Error(t, err)
	assert.Nil(t, created)
	mockRepo.AssertExpectations(t)
}

func TestRegister_FailToSaveUser(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)
	user := &userv1.User{
		Name:     "John Johnson",
		Barcode:  "123456",
		Password: "12345678",
		Role:     "student",
	}
	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(nil, nil)
	mockRepo.On("SaveUser", mock.Anything, mock.AnythingOfType("*userv1.User")).Return(status.Error(codes.Internal, errs.ErrOnUserSaving.Error()))
	created, err := svc.Register(context.Background(), user.Name, user.Barcode, user.Password, user.Role)
	assert.Error(t, err)
	assert.Nil(t, created)
	mockRepo.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)
	user := &userv1.User{
		Name:     "John Johnson",
		Barcode:  "123456",
		Password: "12345678",
		Role:     "student",
	}
	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").Return(&userv1.User{}, nil)
	created, err := svc.Register(context.Background(), user.Name, user.Barcode, user.Password, user.Role)
	assert.Error(t, err)
	assert.Nil(t, created)
	mockRepo.AssertExpectations(t)

}

func TestRegister_ShortPassword(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)

	created, err := svc.Register(context.Background(), "John", "123456", "123", "student")
	assert.Error(t, err)
	assert.Nil(t, created)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_GetUserByBarcodeError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	limiter := &fakeLimiter{}
	cfg := config.MustLoad()
	svc := service.NewUserService(mockRepo, limiter, cfg, nil)

	mockRepo.On("GetUserByBarcode", mock.Anything, "123456").
		Return(nil, status.Error(codes.Internal, "failed to check user"))

	created, err := svc.Register(context.Background(), "John", "123456", "12345678", "student")
	assert.Error(t, err)
	assert.Nil(t, created)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
