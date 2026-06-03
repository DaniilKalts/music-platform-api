package admin

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type TrackRepository interface {
	CreateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	UpdateTrack(ctx context.Context, t *track.Track) (*track.Track, error)
	SoftDeleteTrack(ctx context.Context, id uuid.UUID) error
	GetTrackByID(ctx context.Context, id uuid.UUID) (*track.Track, error)

	FindOrCreateArtist(ctx context.Context, a *track.Artist) (*track.Artist, error)
	FindOrCreateAlbum(ctx context.Context, a *track.Album) (*track.Album, error)
	GetGenreByID(ctx context.Context, id uuid.UUID) (*track.Genre, error)
}

type UserRepository interface {
	UpdateSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error)
}

type TrackCache interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	trackRepo  TrackRepository
	userRepo   UserRepository
	trackCache TrackCache
}

func NewService(trackRepo TrackRepository, userRepo UserRepository, trackCache TrackCache) *Service {
	return &Service{
		trackRepo:  trackRepo,
		userRepo:   userRepo,
		trackCache: trackCache,
	}
}

type CreateTrackInput struct {
	Title           string
	ArtistName      string
	AlbumName       string
	GenreID         uuid.UUID
	DurationSeconds int
	FileURL         string
}

func (s *Service) CreateTrack(ctx context.Context, input CreateTrackInput) (*track.Track, error) {
	// 1. Find or create artist
	artist, err := track.NewArtist(input.ArtistName)
	if err != nil {
		return nil, err
	}
	artist, err = s.trackRepo.FindOrCreateArtist(ctx, artist)
	if err != nil {
		return nil, err
	}

	// 2. Find or create album
	album, err := track.NewAlbum(input.AlbumName)
	if err != nil {
		return nil, err
	}
	album, err = s.trackRepo.FindOrCreateAlbum(ctx, album)
	if err != nil {
		return nil, err
	}

	// 3. Check genre
	_, err = s.trackRepo.GetGenreByID(ctx, input.GenreID)
	if err != nil {
		return nil, err
	}

	// 4. Create track
	t, err := track.NewTrack(input.Title, artist.ID, album.ID, input.GenreID, input.DurationSeconds, input.FileURL)
	if err != nil {
		return nil, err
	}

	return s.trackRepo.CreateTrack(ctx, t)
}

type UpdateTrackInput struct {
	ID              uuid.UUID
	Title           string
	ArtistName      string
	AlbumName       string
	GenreID         uuid.UUID
	DurationSeconds int
	FileURL         string
}

func (s *Service) UpdateTrack(ctx context.Context, input UpdateTrackInput) (*track.Track, error) {
	// 1. Get existing track
	t, err := s.trackRepo.GetTrackByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// 2. Find or create artist
	artist, err := track.NewArtist(input.ArtistName)
	if err != nil {
		return nil, err
	}
	artist, err = s.trackRepo.FindOrCreateArtist(ctx, artist)
	if err != nil {
		return nil, err
	}

	// 3. Find or create album
	album, err := track.NewAlbum(input.AlbumName)
	if err != nil {
		return nil, err
	}
	album, err = s.trackRepo.FindOrCreateAlbum(ctx, album)
	if err != nil {
		return nil, err
	}

	// 4. Check genre
	_, err = s.trackRepo.GetGenreByID(ctx, input.GenreID)
	if err != nil {
		return nil, err
	}

	// 5. Update track
	if err := t.Update(input.Title, artist.ID, album.ID, input.GenreID, input.DurationSeconds, input.FileURL); err != nil {
		return nil, err
	}

	updated, err := s.trackRepo.UpdateTrack(ctx, t)
	if err != nil {
		return nil, err
	}

	// 6. Invalidate cache
	_ = s.trackCache.Delete(ctx, input.ID)

	return updated, nil
}

func (s *Service) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	if err := s.trackRepo.SoftDeleteTrack(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.trackCache.Delete(ctx, id)

	return nil
}

func (s *Service) UpdateUserSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error) {
	return s.userRepo.UpdateSubscription(ctx, id, sub)
}
