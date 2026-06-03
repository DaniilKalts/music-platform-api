package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"

	"github.com/DaniilKalts/music-platform-api/internal/repository/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/repository/history"
	"github.com/DaniilKalts/music-platform-api/internal/repository/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/repository/track"
	"github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, email, username *string) (*user.User, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error)
}

type TrackRepository interface {
	CreateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*track.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*track.Track, error)
	UpdateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	SoftDeleteTrack(ctx context.Context, id uuid.UUID) error
	TrackExists(ctx context.Context, id uuid.UUID) (bool, error)

	CreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error)
	GetArtistByID(ctx context.Context, id uuid.UUID) (*track.Artist, error)
	FindOrCreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error)

	CreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error)
	GetAlbumByID(ctx context.Context, id uuid.UUID) (*track.Album, error)
	FindOrCreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error)

	CreateGenre(ctx context.Context, g *track.Genre) (*track.Genre, error)
	ListGenres(ctx context.Context) ([]*track.Genre, error)
	GetGenreByID(ctx context.Context, id uuid.UUID) (*track.Genre, error)

	CreateTrackWithDependencies(ctx context.Context, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
	UpdateTrackWithDependencies(ctx context.Context, id uuid.UUID, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
}

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error)
	ListPlaylistsByUserID(ctx context.Context, userID uuid.UUID) ([]*playlist.Playlist, error)
	GetPlaylistByIDForUser(ctx context.Context, id, userID uuid.UUID) (*playlist.Playlist, error)
	UpdatePlaylist(ctx context.Context, p *playlist.Playlist) (*playlist.Playlist, error)
	DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error
	CountPlaylistsByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	AddTrackToPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	ListPlaylistTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*track.Track, error)
}

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, f *favorite.Favorite) error
	RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	ListFavoritesByUserID(ctx context.Context, userID uuid.UUID) ([]*track.Track, error)
	CountFavoritesByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	FavoriteExists(ctx context.Context, userID, trackID uuid.UUID) (bool, error)
}

type HistoryRepository interface {
	CreateListeningHistory(ctx context.Context, h *history.HistoryRecord) error
	ListListeningHistoryByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error)
}

type Repositories struct {
	User     UserRepository
	Track    TrackRepository
	Playlist PlaylistRepository
	Favorite FavoriteRepository
	History  HistoryRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:     userrepo.NewRepository(db),
		Track:    trackrepo.NewRepository(db),
		Playlist: playlistrepo.NewRepository(db),
		Favorite: favoriterepo.NewRepository(db),
		History:  historyrepo.NewRepository(db),
	}
}
