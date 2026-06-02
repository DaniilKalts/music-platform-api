package jwt

import "errors"

var ErrInvalidToken = errors.New("invalid or expired token")
