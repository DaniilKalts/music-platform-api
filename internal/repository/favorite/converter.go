package favoriterepo

import (
	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

func toDomainTrackFromFavoriteList(row sqlc.ListFavoritesByUserIDRow) *track.Track {
	return &track.Track{
		ID:              row.ID,
		Title:           row.Title,
		ArtistID:        row.ArtistID,
		AlbumID:         row.AlbumID,
		GenreID:         row.GenreID,
		DurationSeconds: int(row.DurationSeconds),
		FileURL:         row.FileUrl,
		DeletedAt:       row.DeletedAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		ArtistName:      row.ArtistName,
		AlbumName:       row.AlbumName,
		GenreName:       row.GenreName,
	}
}
