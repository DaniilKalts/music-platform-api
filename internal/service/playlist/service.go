package playlist

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

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

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type Service struct {
	playlists PlaylistRepository
	users     UserRepository
	freeLimit int
}

func NewService(playlists PlaylistRepository, users UserRepository, freeLimit int) *Service {
	return &Service{
		playlists: playlists,
		users:     users,
		freeLimit: freeLimit,
	}
}

type CreateInput struct {
	UserID      uuid.UUID
	Name        string
	Description *string
}

func (s *Service) CreatePlaylist(ctx context.Context, input CreateInput) (*playlist.Playlist, error) {
	u, err := s.users.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if u.Subscription == user.SubscriptionFree {
		count, err := s.playlists.CountPlaylistsByUserID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("count playlists: %w", err)
		}

		if int(count) >= s.freeLimit {
			return nil, playlist.ErrPlaylistLimitReached
		}
	}

	p, err := playlist.NewPlaylist(input.UserID, input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	return s.playlists.CreatePlaylist(ctx, p)
}

func (s *Service) GetPlaylist(ctx context.Context, id, userID uuid.UUID) (*playlist.Playlist, error) {
	return s.playlists.GetPlaylistByIDForUser(ctx, id, userID)
}

func (s *Service) ListPlaylists(ctx context.Context, userID uuid.UUID) ([]*playlist.Playlist, error) {
	return s.playlists.ListPlaylistsByUserID(ctx, userID)
}

type UpdateInput struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description *string
}

func (s *Service) UpdatePlaylist(ctx context.Context, input UpdateInput) (*playlist.Playlist, error) {
	p, err := s.playlists.GetPlaylistByIDForUser(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	if err := p.Update(input.Name, input.Description); err != nil {
		return nil, err
	}

	return s.playlists.UpdatePlaylist(ctx, p)
}

func (s *Service) DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error {
	return s.playlists.DeletePlaylist(ctx, id, userID)
}

func (s *Service) AddTrack(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	return s.playlists.AddTrackToPlaylist(ctx, playlistID, trackID, userID)
}

func (s *Service) RemoveTrack(ctx context.Context, playlistID, trackID, userID uuid.UUID) error {
	return s.playlists.RemoveTrackFromPlaylist(ctx, playlistID, trackID, userID)
}

func (s *Service) ListTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*track.Track, error) {
	return s.playlists.ListPlaylistTracks(ctx, playlistID, userID)
}
