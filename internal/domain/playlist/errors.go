package playlist

import "errors"

var (
	ErrPlaylistNotFound     = errors.New("playlist not found")
	ErrPlaylistLimitReached = errors.New("playlist limit reached")
	ErrInvalidPlaylistName  = errors.New("invalid playlist name")
	ErrInvalidDescription   = errors.New("invalid playlist description")
	ErrOwnerRequired        = errors.New("playlist owner is required")
)

var fieldErrors = map[string]error{
	"Name":        ErrInvalidPlaylistName,
	"Description": ErrInvalidDescription,
	"UserID":      ErrOwnerRequired,
}
