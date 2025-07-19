package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")
	ErrUuidIsEmpty  = errors.New("uuid cannot be empty")
)
