package track

import "errors"

var (
	ErrTrackNotFound  = errors.New("track not found")
	ErrArtistNotFound = errors.New("artist not found")
	ErrAlbumNotFound  = errors.New("album not found")
	ErrGenreNotFound  = errors.New("genre not found")

	ErrInvalidTitle     = errors.New("invalid track title")
	ErrInvalidDuration  = errors.New("invalid track duration")
	ErrInvalidFileURL   = errors.New("invalid file url")
	ErrInvalidName      = errors.New("invalid name (artist/album/genre)")
	ErrArtistIDRequired = errors.New("artist id is required")
	ErrAlbumIDRequired  = errors.New("album id is required")
	ErrGenreIDRequired  = errors.New("genre id is required")
)

var fieldErrors = map[string]error{
	"Title":           ErrInvalidTitle,
	"DurationSeconds": ErrInvalidDuration,
	"FileURL":         ErrInvalidFileURL,
	"Name":            ErrInvalidName,
	"ArtistID":        ErrArtistIDRequired,
	"AlbumID":         ErrAlbumIDRequired,
	"GenreID":         ErrGenreIDRequired,
}
