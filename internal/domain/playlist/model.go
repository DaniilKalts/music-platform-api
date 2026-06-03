package playlist

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID `validate:"required"`
	UserID      uuid.UUID `validate:"required"`
	Name        string    `validate:"required,min=1,max=100"`
	Description *string   `validate:"omitempty,max=500"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPlaylist(userID uuid.UUID, name string, description *string) (*Playlist, error) {
	p := &Playlist{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        strings.TrimSpace(name),
		Description: description,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Playlist) Update(name string, description *string) error {
	p.Name = strings.TrimSpace(name)
	p.Description = description

	return p.Validate()
}

func (p *Playlist) IsOwner(userID uuid.UUID) bool {
	return p.UserID == userID
}

func (p *Playlist) Validate() error {
	return mapValidationError(validate.Struct(p))
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
