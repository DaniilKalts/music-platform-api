package history

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type HistoryRecord struct {
	ID         uuid.UUID `validate:"required"`
	UserID     uuid.UUID `validate:"required"`
	TrackID    uuid.UUID `validate:"required"`
	ListenedAt time.Time
}

func NewHistoryRecord(userID, trackID uuid.UUID) (*HistoryRecord, error) {
	h := &HistoryRecord{
		ID:      uuid.New(),
		UserID:  userID,
		TrackID: trackID,
	}

	if err := h.Validate(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *HistoryRecord) Validate() error {
	return mapValidationError(validate.Struct(h))
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
