package model

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserPassword    = errors.New("user password is wrong")
	ErrSessionCreate   = errors.New("session create error")
	ErrSessionNotFound = errors.New("session not found")
)
