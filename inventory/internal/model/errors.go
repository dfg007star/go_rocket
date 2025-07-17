package model

import "errors"

var ErrPartNotFound = errors.New("part not found")
var ErrUuidIsEmpty = errors.New("uuid cannot be empty")
