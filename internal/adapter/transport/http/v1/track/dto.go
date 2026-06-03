package track

import (
	"time"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

type TrackResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	ArtistID        string    `json:"artist_id"`
	ArtistName      string    `json:"artist_name"`
	AlbumID         string    `json:"album_id"`
	AlbumName       string    `json:"album_name"`
	GenreID         string    `json:"genre_id"`
	GenreName       string    `json:"genre_name"`
	DurationSeconds int       `json:"duration_seconds"`
	FileURL         string    `json:"file_url"`
	CreatedAt       time.Time `json:"created_at"`
}

func ToTrackResponse(t track.Track) TrackResponse {
	return TrackResponse{
		ID:              t.ID.String(),
		Title:           t.Title,
		ArtistID:        t.ArtistID.String(),
		ArtistName:      t.ArtistName,
		AlbumID:         t.AlbumID.String(),
		AlbumName:       t.AlbumName,
		GenreID:         t.GenreID.String(),
		GenreName:       t.GenreName,
		DurationSeconds: t.DurationSeconds,
		FileURL:         t.FileURL,
		CreatedAt:       t.CreatedAt,
	}
}

func ToTracksResponse(tracks []*track.Track) []TrackResponse {
	res := make([]TrackResponse, len(tracks))
	for i, t := range tracks {
		res[i] = ToTrackResponse(*t)
	}
	return res
}

type GenreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToGenresResponse(genres []*track.Genre) []GenreResponse {
	res := make([]GenreResponse, len(genres))
	for i, g := range genres {
		res[i] = GenreResponse{
			ID:   g.ID.String(),
			Name: g.Name,
		}
	}
	return res
}
