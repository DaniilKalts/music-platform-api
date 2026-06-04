package playlist

import (
	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
)

type CreatePlaylistRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

type UpdatePlaylistRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

type PlaylistResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func ToPlaylistResponse(p playlist.Playlist) PlaylistResponse {
	return PlaylistResponse{
		ID:          p.ID.String(),
		UserID:      p.UserID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToPlaylistsResponse(playlists []*playlist.Playlist) []PlaylistResponse {
	res := make([]PlaylistResponse, len(playlists))
	for i, p := range playlists {
		res[i] = ToPlaylistResponse(*p)
	}
	return res
}

type AddTrackRequest struct {
	TrackID string `json:"track_id" validate:"required,uuid"`
}
