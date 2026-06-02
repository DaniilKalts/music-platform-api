package httpx

import (
	"errors"
	"net/http"
	"strings"
)

var ErrEmptyAuthHeader = errors.New("empty auth header")
var ErrInvalidAuthHeader = errors.New("invalid auth header")

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
