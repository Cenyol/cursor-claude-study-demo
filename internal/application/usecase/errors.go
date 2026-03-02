package usecase

import "errors"

var (
	ErrPasswordTooShort = errors.New("password too short")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUsernameExists   = errors.New("username already exists")
)
