package errors

import "errors"

var (
	ErrWSInvalidFormat = errors.New("invalid web socket message format ")
	ErrCannotParseClaims = errors.New("cannot parse claims")
	ErrJWTIsExpired = errors.New("jwt is expired")
)
