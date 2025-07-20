package model

import "errors"

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrNotAllPartsMatched    = errors.New("not all orders matched")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderInternalError    = errors.New("internal error")
)
