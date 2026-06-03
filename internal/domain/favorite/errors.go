package favorite

import "errors"

var (
	ErrFavoriteAlreadyExists = errors.New("track already in favorites")
	ErrFavoriteLimitReached  = errors.New("favorites limit reached")
	ErrUserIDRequired        = errors.New("user id is required")
	ErrTrackIDRequired       = errors.New("track id is required")
)

var fieldErrors = map[string]error{
	"UserID":  ErrUserIDRequired,
	"TrackID": ErrTrackIDRequired,
}
