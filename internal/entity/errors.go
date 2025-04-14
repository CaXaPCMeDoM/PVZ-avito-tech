package entity

import "errors"

var (
	ErrPasswordTooLong   = errors.New("password too long")
	ErrPasswordHashing   = errors.New("failed to process password")
	ErrInvalidPassword   = errors.New("invalid password format")
	ErrPasswordVerify    = errors.New("passwords don't match")
	ErrUserNotFound      = errors.New("user not found")
	ErrInternal          = errors.New("internal error")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrCreatePVZ  = errors.New("failed to create PVZ")
	ErrGetPVZList = errors.New("failed to get PVZ list")

	ErrPVZNotFound       = errors.New("pvz not found")
	ErrReceptionConflict = errors.New("existing open reception")

	ErrNoActiveReception = errors.New("no active reception")
	ErrNoProducts        = errors.New("no products")
)
