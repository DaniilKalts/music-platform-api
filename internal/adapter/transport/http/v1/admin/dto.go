package admin

import (
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type CreateTrackRequest struct {
	Title           string `json:"title" validate:"required,min=1,max=255"`
	ArtistName      string `json:"artist_name" validate:"required,min=1,max=255"`
	AlbumName       string `json:"album_name" validate:"required,min=1,max=255"`
	GenreID         string `json:"genre_id" validate:"required,uuid"`
	DurationSeconds int    `json:"duration_seconds" validate:"required,min=1"`
	FileURL         string `json:"file_url" validate:"required,url"`
}

type UpdateTrackRequest struct {
	Title           string `json:"title" validate:"required,min=1,max=255"`
	ArtistName      string `json:"artist_name" validate:"required,min=1,max=255"`
	AlbumName       string `json:"album_name" validate:"required,min=1,max=255"`
	GenreID         string `json:"genre_id" validate:"required,uuid"`
	DurationSeconds int    `json:"duration_seconds" validate:"required,min=1"`
	FileURL         string `json:"file_url" validate:"required,url"`
}

type UpdateSubscriptionRequest struct {
	Type user.Subscription `json:"type" validate:"required,oneof=FREE PREMIUM"`
}
