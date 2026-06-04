package user

import "errors"

var (
	ErrEmailInvalid     = errors.New("email invalid")
	ErrPasswordTooShort = errors.New("password too short")
	ErrInvalidRole      = errors.New("invalid role")
	ErrUserNotFound     = errors.New("user not found")
)
