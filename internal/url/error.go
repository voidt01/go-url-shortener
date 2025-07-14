package url

import (
	"errors"
)

var (
	ErrUrlNotFound = errors.New("models: No matching url found")
	ErrUrlInvalid  = errors.New("models: Url format is invalid")
)
