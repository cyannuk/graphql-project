package core

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrNoValue          = errors.New("no value")
	ErrInvalidValueType = errors.New("invalid value value")
)
