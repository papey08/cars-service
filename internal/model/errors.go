package model

import "errors"

var (
	ErrDuplicateRegNum = errors.New("duplicate registration number")
	ErrInvalidInput    = errors.New("invalid input")
	ErrValidation      = errors.New("validation error")
	ErrCarNotFound     = errors.New("car not found")
	ErrDatabaseError   = errors.New("database error")
	ErrApiError        = errors.New("outer api error")
	ErrServiceError    = errors.New("unknown service error")
)
