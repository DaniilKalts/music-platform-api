package admin

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

type TrackRepository interface {
	CreateTrackWithDependencies(ctx context.Context, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
	UpdateTrackWithDependencies(ctx context.Context, id uuid.UUID, title, artistName, albumName string, genreID uuid.UUID, durationSeconds int, fileURL string) (*track.Track, error)
	SoftDeleteTrack(ctx context.Context, id uuid.UUID) error
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
	return s.trackRepo.CreateTrackWithDependencies(ctx, input.Title, input.ArtistName, input.AlbumName, input.GenreID, input.DurationSeconds, input.FileURL)
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
	updated, err := s.trackRepo.UpdateTrackWithDependencies(ctx, input.ID, input.Title, input.ArtistName, input.AlbumName, input.GenreID, input.DurationSeconds, input.FileURL)
	if err != nil {
		return nil, err
	}

	_ = s.trackCache.Delete(ctx, input.ID)

	return updated, nil
}

func (s *Service) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	if err := s.trackRepo.SoftDeleteTrack(ctx, id); err != nil {
		return err
	}

	_ = s.trackCache.Delete(ctx, id)

	return nil
}

func (s *Service) UpdateUserSubscription(ctx context.Context, id uuid.UUID, sub user.Subscription) (*user.User, error) {
	return s.userRepo.UpdateSubscription(ctx, id, sub)
}
