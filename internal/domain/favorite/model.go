package favorite

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Favorite struct {
	UserID    uuid.UUID `validate:"required"`
	TrackID   uuid.UUID `validate:"required"`
	CreatedAt time.Time
}

func NewFavorite(userID, trackID uuid.UUID) (*Favorite, error) {
	f := &Favorite{
		UserID:  userID,
		TrackID: trackID,
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Favorite) Validate() error {
	return mapValidationError(validate.Struct(f))
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
