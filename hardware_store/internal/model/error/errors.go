package model

import "errors"

var ErrImageNotFound = errors.New("image not found")

var ErrProductNotFound error = errors.New("product not found")
var ErrInsufficientStock = errors.New("insufficient stock")
var ErrAmountIsNegative = errors.New("amount must be positive")