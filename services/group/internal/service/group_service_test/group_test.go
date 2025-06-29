package group_service_test

import (
	"context"
	"errors"
	groupv1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/group/internal/config"
	"github.com/mdqni/Attendly/services/group/internal/repository/mocks"
	"github.com/mdqni/Attendly/services/group/internal/service"
	"github.com/mdqni/Attendly/shared/errs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGroupService_CreateGroup_Success(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})

	ctx := context.Background()

	mockRepo.On("CreateGroup", ctx, "SE-2430", "SE", 2024).
		Return(&groupv1.Group{Id: "123"}, nil)

	group, err := svc.CreateGroup(ctx, "SE-2430", "SE", 2024)

	require.NoError(t, err)
	require.Equal(t, "123", group.Id)
	mockRepo.AssertCalled(t, "CreateGroup", ctx, "SE-2430", "SE", 2024)
}

func TestGroupService_CreateGroup_GroupAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	mockRepo.On("CreateGroup", ctx, "SE-2430", "SE", 2024).Return(nil,
		errs.ErrGroupAlreadyExists)
	group, err := svc.CreateGroup(ctx, "SE-2430", "SE", 2024)
	require.Nil(t, group)
	require.Error(t, err)
	require.Equal(t, codes.AlreadyExists, status.Code(err))

	mockRepo.AssertCalled(t, "CreateGroup", ctx, "SE-2430", "SE", 2024)
}

func TestGroupService_CreateGroup_InternalError(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	mockRepo.On("CreateGroup", ctx, "SE-2430", "SE", 2024).Return(nil, status.Error(codes.Internal, "internal error"))
	group, err := svc.CreateGroup(ctx, "SE-2430", "SE", 2024)
	require.Nil(t, group)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
	mockRepo.AssertCalled(t, "CreateGroup", ctx, mock.Anything, mock.Anything, mock.Anything)
}

func TestGroupService_GetGroup_Success(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	mockRepo.On("GetGroup", ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed1").Return(
		&groupv1.Group{
			Id:         "a9f3d25d-05b3-4f04-94bc-78405513eed1",
			Name:       "SE-2433",
			Department: "SE",
			Year:       2024,
		}, nil)
	group, err := svc.GetGroup(ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed1")
	require.NoError(t, err)
	require.Equal(t, "SE-2433", group.Name)
	mockRepo.AssertCalled(t, "GetGroup", ctx, mock.Anything, mock.Anything)
}
func TestGroupService_GetGroup_NotFound(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	mockRepo.On("GetGroup", ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed").Return(
		nil, errs.ErrGroupNotFound)
	group, err := svc.GetGroup(ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed")
	require.Nil(t, group)
	require.Error(t, err)
	require.Equal(t, status.Error(codes.NotFound, "group not found"), err)
	mockRepo.AssertCalled(t, "GetGroup", ctx, mock.Anything, mock.Anything)
}

func TestGroupService_GetGroup_InternalError(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	mockRepo.On("GetGroup", ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed").Return(
		nil, errors.New("db is down"))
	group, err := svc.GetGroup(ctx, "a9f3d25d-05b3-4f04-94bc-78405513eed")

	require.Nil(t, group)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, st.Code())
	require.Equal(t, "Internal error", st.Message())
	mockRepo.AssertCalled(t, "GetGroup", ctx, mock.Anything, mock.Anything)
}

func TestGroupService_AddUserToGroup_Success(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"

	mockRepo.On("AddUserToGroup", ctx, groupID, userID).Return(true, nil)

	ok, err := svc.AddUserToGroup(ctx, groupID, userID)

	require.NoError(t, err)
	require.True(t, ok)
	mockRepo.AssertCalled(t, "AddUserToGroup", ctx, groupID, userID)
}

func TestGroupService_AddUserToGroup_UserOrGroupNotExists(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"

	mockRepo.On("AddUserToGroup", ctx, groupID, userID).Return(false, errs.ErrUserOrGroupNotExists)

	ok, err := svc.AddUserToGroup(ctx, groupID, userID)
	require.False(t, ok)
	require.Error(t, err)
	require.Equal(t, status.Error(codes.NotFound, "user or group not exists"), err)
	mockRepo.AssertCalled(t, "AddUserToGroup", ctx, groupID, userID)
}

func TestGroupService_AddUserToGroup_InternalError(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"
	mockRepo.On("AddUserToGroup", ctx, groupID, userID).Return(false, errors.New("error adding user to group"))
	ok, err := svc.AddUserToGroup(ctx, groupID, userID)
	require.False(t, ok)
	require.Error(t, err)
	require.Equal(t, status.Error(codes.Internal, "error adding user to group"), err)
	mockRepo.AssertCalled(t, "AddUserToGroup", ctx, groupID, userID)
}

func TestGroupService_RemoveUserFromGroup_Success(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"
	mockRepo.On("RemoveUserFromGroup", ctx, groupID, userID).Return(true, nil)
	ok, err := svc.RemoveUserFromGroup(ctx, groupID, userID)
	require.NoError(t, err)
	require.True(t, ok)
	mockRepo.AssertCalled(t, "RemoveUserFromGroup", ctx, groupID, userID)
}
func TestGroupService_RemoveUserFromGroup_UserOrGroupNotExists(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"
	mockRepo.On("RemoveUserFromGroup", ctx, groupID, userID).Return(false, errs.ErrUserOrGroupNotExists)
	ok, err := svc.RemoveUserFromGroup(ctx, groupID, userID)
	require.False(t, ok)
	require.Error(t, err)
	require.Equal(t, status.Error(codes.NotFound, "user or group not exists"), err)
	mockRepo.AssertCalled(t, "RemoveUserFromGroup", ctx, groupID, userID)
}
func TestGroupService_RemoveUserFromGroup_InternalError(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	userID := "u-456"
	mockRepo.On("RemoveUserFromGroup", ctx, groupID, userID).Return(false, errors.New("error removing user from group"))
	ok, err := svc.RemoveUserFromGroup(ctx, groupID, userID)
	require.False(t, ok)
	require.Error(t, err)
	require.Equal(t, status.Error(codes.Internal, "Internal error"), err)
	mockRepo.AssertCalled(t, "RemoveUserFromGroup", ctx, groupID, userID)
}

func TestGroupService_ListUsersInGroup_Success(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	mockRepo.On("ListUsersInGroup", ctx, groupID).Return([]*userv1.User{}, nil)
	users, err := svc.ListUsersInGroup(ctx, groupID)
	require.NoError(t, err)
	require.Equal(t, []*userv1.User{}, users)
	mockRepo.AssertCalled(t, "ListUsersInGroup", ctx, groupID)
}

func TestGroupService_ListUsersInGroup_Empty(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"

	mockRepo.On("ListUsersInGroup", ctx, groupID).Return(nil, errs.ErrGroupIsEmpty)

	users, err := svc.ListUsersInGroup(ctx, groupID)

	require.Nil(t, users)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Equal(t, "group is empty", st.Message())

	mockRepo.AssertCalled(t, "ListUsersInGroup", ctx, groupID)
}

func TestGroupService_ListUsersInGroup_InternalError(t *testing.T) {
	mockRepo := mocks.NewGroupRepository(t)
	svc := service.NewGroupService(mockRepo, &config.Config{})
	ctx := context.Background()
	groupID := "g-123"
	mockRepo.On("ListUsersInGroup", ctx, groupID).Return(nil, errors.New("internal error"))

	users, err := svc.ListUsersInGroup(ctx, groupID)

	require.Nil(t, users)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Internal, st.Code())
	require.Equal(t, "internal error", st.Message())

	mockRepo.AssertCalled(t, "ListUsersInGroup", ctx, groupID)

}
