package model

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrValidation    = errors.New("validation error")
	ErrCarNotFound   = errors.New("car not found")
	ErrDatabaseError = errors.New("database error")
	ErrApiError      = errors.New("outer api error")
)
