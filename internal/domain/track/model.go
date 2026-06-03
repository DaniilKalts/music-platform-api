package track

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Track struct {
	ID              uuid.UUID  `validate:"required"`
	Title           string     `validate:"required,min=1,max=255"`
	ArtistID        uuid.UUID  `validate:"required"`
	AlbumID         uuid.UUID  `validate:"required"`
	GenreID         uuid.UUID  `validate:"required"`
	DurationSeconds int        `validate:"required,gt=0"`
	FileURL         string     `validate:"required,url"`
	DeletedAt       *time.Time `validate:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	ArtistName string `validate:"-"`
	AlbumName  string `validate:"-"`
	GenreName  string `validate:"-"`
}

func NewTrack(title string, artistID, albumID, genreID uuid.UUID, duration int, fileURL string) (*Track, error) {
	t := &Track{
		ID:              uuid.New(),
		Title:           strings.TrimSpace(title),
		ArtistID:        artistID,
		AlbumID:         albumID,
		GenreID:         genreID,
		DurationSeconds: duration,
		FileURL:         strings.TrimSpace(fileURL),
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Track) Update(title string, artistID, albumID, genreID uuid.UUID, duration int, fileURL string) error {
	t.Title = strings.TrimSpace(title)
	t.ArtistID = artistID
	t.AlbumID = albumID
	t.GenreID = genreID
	t.DurationSeconds = duration
	t.FileURL = strings.TrimSpace(fileURL)

	return t.Validate()
}

func (t *Track) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *Track) Validate() error {
	return mapValidationError(validate.Struct(t))
}

type Artist struct {
	ID        uuid.UUID `validate:"required"`
	Name      string    `validate:"required,min=1,max=255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewArtist(name string) (*Artist, error) {
	a := &Artist{
		ID:   uuid.New(),
		Name: strings.TrimSpace(name),
	}
	if err := mapValidationError(validate.Struct(a)); err != nil {
		return nil, err
	}
	return a, nil
}

type Album struct {
	ID        uuid.UUID `validate:"required"`
	Name      string    `validate:"required,min=1,max=255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAlbum(name string) (*Album, error) {
	a := &Album{
		ID:   uuid.New(),
		Name: strings.TrimSpace(name),
	}
	if err := mapValidationError(validate.Struct(a)); err != nil {
		return nil, err
	}
	return a, nil
}

type Genre struct {
	ID        uuid.UUID `validate:"required"`
	Name      string    `validate:"required,min=1,max=100"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
