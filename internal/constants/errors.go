package constants

import "errors"

var (
	ErrInvalidAuthHeader    = errors.New("missing or invalid authorization header")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrExpiredJWT           = errors.New("JWT expired")
)
