package user

import "errors"

var (
	ErrNotFound              = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("user with this email already exists")
	ErrUsernameAlreadyExists = errors.New("user with this username already exists")
	ErrInvalidEmail          = errors.New("invalid user email")
	ErrInvalidUsername       = errors.New("invalid username")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidRole           = errors.New("invalid user role")
	ErrInvalidSubscription   = errors.New("invalid user subscription")
)

var fieldErrors = map[string]error{
	"Email":        ErrInvalidEmail,
	"Username":     ErrInvalidUsername,
	"Role":         ErrInvalidRole,
	"Subscription": ErrInvalidSubscription,
}
