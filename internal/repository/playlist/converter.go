package playlistrepo

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

func toDomainPlaylistFromCreate(row sqlc.Playlist) *playlist.Playlist {
	return &playlist.Playlist{
		ID:          row.ID,
		UserID:      row.UserID,
		Name:        row.Name,
		Description: fromPgText(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func toDomainPlaylistFromList(row sqlc.Playlist) *playlist.Playlist {
	return &playlist.Playlist{
		ID:          row.ID,
		UserID:      row.UserID,
		Name:        row.Name,
		Description: fromPgText(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func toDomainPlaylistFromGet(row sqlc.Playlist) *playlist.Playlist {
	return &playlist.Playlist{
		ID:          row.ID,
		UserID:      row.UserID,
		Name:        row.Name,
		Description: fromPgText(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func toDomainPlaylistFromUpdate(row sqlc.Playlist) *playlist.Playlist {
	return &playlist.Playlist{
		ID:          row.ID,
		UserID:      row.UserID,
		Name:        row.Name,
		Description: fromPgText(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func toDomainTrackFromPlaylistList(row sqlc.ListPlaylistTracksRow) *track.Track {
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

func fromPgText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
