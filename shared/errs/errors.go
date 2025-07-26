package errs

import "errors"

var (
	ErrNoTopics        = errors.New("no topics provided")
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token expired")
	ErrInvalidPassword = errors.New("invalid password")
	//-------------------------
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrMissingField         = errors.New("missing required fields")
	ErrTooManyLoginAttempt  = errors.New("too many login attempts, try again later")
	ErrPasswordTooShort     = errors.New("password is too short")
	ErrOnUserSaving         = errors.New("failed to register user")
	ErrUserOrGroupNotExists = errors.New("user or group not exists")
	/*-------------------------------------------------------------------*/
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrNoUserInGroup      = errors.New("no such user in group")
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupIsEmpty       = errors.New("group is empty")
)
