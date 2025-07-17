package url

import (
	"errors"
)

var (
	ErrUrlNotFound   = errors.New("store: No matching url found")
	ErrUrlInvalid    = errors.New("store: Url format is invalid")
	ErrUrlDuplicated = errors.New("store: There's already existed url")
	ErrShortUrlFailedGeneration = errors.New("service: Failed to generate shortCode")

	ErrDuplicateEntryCode uint16 = 1062
)
