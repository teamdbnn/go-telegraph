package telegraph

import (
	"github.com/pkg/errors"
)

var (
	// ErrInvalidDataType is returned when ContentFormat function are passed a data argument of invalid type.
	ErrInvalidDataType = errors.New("invalid data type")

	// ErrNoInputData is returned when any method get nil argument.
	ErrNoInputData = errors.New("no input data")

	ErrEmptyAccessToken = errors.New("empty access_token")
)
