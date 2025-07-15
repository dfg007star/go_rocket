package model

import "errors"

var ErrPaymentNotFound = errors.New("payment not found")
var ErrPayOrderModelValidation = errors.New("invalid order model")
