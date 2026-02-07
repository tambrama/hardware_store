package storage

import (
	"errors"
)

var (
	ErrClientNotFound    = errors.New("client not found")
	ErrImageNotFound     = errors.New("image not found")
	ErrAddressNotFound   = errors.New("address not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrSupplierNotFound  = errors.New("supplier not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrClientExists      = errors.New("client exists")
	ErrCreation          = errors.New("сreation error")
	ErrDelete            = errors.New("delete error")
	ErrUpdate            = errors.New("update error")
)
