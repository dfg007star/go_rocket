package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrNotAllPartsMatched = errors.New("not all orders matched")
var ErrOrderAlreadyPaid = errors.New("order already paid")
var ErrOrderAlreadyCancelled = errors.New("order already cancelled")
