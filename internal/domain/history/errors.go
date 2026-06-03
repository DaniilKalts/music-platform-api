package history

import "errors"

var (
	ErrUserIDRequired  = errors.New("user id is required")
	ErrTrackIDRequired = errors.New("track id is required")
)

var fieldErrors = map[string]error{
	"UserID":  ErrUserIDRequired,
	"TrackID": ErrTrackIDRequired,
}
