package model

import "errors"

var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrPayOrderModelValidation = errors.New("invalid order model")
)
