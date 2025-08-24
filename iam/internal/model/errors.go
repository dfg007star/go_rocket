package model

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrUserPassword              = errors.New("user password is wrong")
	ErrSessionCreate             = errors.New("session create error")
	ErrSessionNotFound           = errors.New("session not found")
	ErrUserLoginRequired         = errors.New("login is required")
	ErrUserEmailRequired         = errors.New("email is required")
	ErrUserPasswordRequired      = errors.New("password is required")
	ErrUserLoginExists           = errors.New("login already exists")
	ErrUserEmailExists           = errors.New("email already exists")
	ErrUserPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrUserEmailInvalid          = errors.New("email format is invalid")
	ErrUserLoginInvalid          = errors.New("login must be between 3 and 50 characters")
	ErrNotificationMethodInvalid = errors.New("notification method is invalid")
)
