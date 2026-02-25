package model

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrRefreshTokenExists   = errors.New("refresh token already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
