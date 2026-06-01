package user

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Email        string       `validate:"required,email,max=254"`
	Username     string       `validate:"required,min=3,max=50"`
	Role         Role         `validate:"required"`
	Subscription Subscription `validate:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, username string) (*User, error) {
	u := &User{
		ID:           uuid.New(),
		Role:         RoleUser,
		Subscription: SubscriptionFree,
	}
	if err := u.UpdateProfile(email, username); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) UpdateProfile(email, username string) error {
	u.Email = NormalizeEmail(email)
	u.Username = strings.TrimSpace(username)

	return u.Validate()
}

func (u *User) Validate() error {
	if err := mapValidationError(validate.Struct(u)); err != nil {
		return err
	}
	if !u.Role.IsValid() {
		return ErrInvalidRole
	}
	if !u.Subscription.IsValid() {
		return ErrInvalidSubscription
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func mapValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	for _, fieldErr := range validationErrors {
		if mapped, ok := fieldErrors[fieldErr.Field()]; ok {
			return mapped
		}
	}

	return err
}
