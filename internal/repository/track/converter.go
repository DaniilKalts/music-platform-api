package trackrepo

import (
	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
)

func toDomainTrackFromCreate(row sqlc.CreateTrackRow) *track.Track {
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

func toDomainTrackFromGet(row sqlc.GetTrackByIDRow) *track.Track {
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

func toDomainTrackFromList(row sqlc.ListTracksRow) *track.Track {
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

func toDomainTrackFromSearch(row sqlc.SearchTracksRow) *track.Track {
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

func toDomainTrackFromUpdate(row sqlc.UpdateTrackRow) *track.Track {
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

func toDomainArtistFromCreate(row sqlc.Artist) *track.Artist {
	return &track.Artist{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainArtistFromGet(row sqlc.Artist) *track.Artist {
	return &track.Artist{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainArtistFromFindOrCreate(row sqlc.Artist) *track.Artist {
	return &track.Artist{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainAlbumFromCreate(row sqlc.Album) *track.Album {
	return &track.Album{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainAlbumFromGet(row sqlc.Album) *track.Album {
	return &track.Album{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainAlbumFromFindOrCreate(row sqlc.Album) *track.Album {
	return &track.Album{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainGenreFromCreate(row sqlc.Genre) *track.Genre {
	return &track.Genre{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainGenreFromGet(row sqlc.Genre) *track.Genre {
	return &track.Genre{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func toDomainGenreFromList(row sqlc.Genre) *track.Genre {
	return &track.Genre{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
