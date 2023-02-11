package httperrors

import (
	"errors"
)

var (
	ErrRecoverPanic      = errors.New("recovering from error")
	ErrInternalServer    = errors.New("something went wrong")
	ErrInvalidPeriod     = errors.New("invalid period")
	ErrInvalidTimezone   = errors.New("invalid timezone")
	ErrInvalidStartPoint = errors.New("invalid start point")
	ErrInvalidEndPoint   = errors.New("invalid end point")
)
