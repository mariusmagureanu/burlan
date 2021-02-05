package errors

import "errors"

var (
	ErrWSInvalidFormat = errors.New("invalid web socket message format ")
)
