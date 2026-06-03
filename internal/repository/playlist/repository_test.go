package playlistrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository/playlist"
	"github.com/DaniilKalts/music-platform-api/internal/repository/testutil"
	trackrepo "github.com/DaniilKalts/music-platform-api/internal/repository/track"
	userrepo "github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

func TestPlaylistRepository(t *testing.T) {
	pool, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	q := sqlc.New(pool)
	repo := playlistrepo.NewRepository(q)
	uRepo := userrepo.NewRepository(pool)
	tRepo := trackrepo.NewRepository(q)
	ctx := context.Background()

	// Helpers
	createUser := func(email string) *user.User {
		u, _ := user.NewUser(email, email)
		p, _ := user.NewPassword("pass")
		created, _ := uRepo.Create(ctx, *u, p)
		return created
	}

	createTrack := func(title string) *track.Track {
		g, _ := tRepo.CreateGenre(ctx, &track.Genre{ID: uuid.New(), Name: title + "G"})
		ar, _ := tRepo.CreateArtist(ctx, &track.Artist{ID: uuid.New(), Name: title + "Ar"})
		al, _ := tRepo.CreateAlbum(ctx, &track.Album{ID: uuid.New(), Name: title + "Al"})
		tr, _ := tRepo.CreateTrack(ctx, &track.Track{ID: uuid.New(), Title: title, ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID})
		return tr
	}

	t.Run("CreateAndGet", func(t *testing.T) {
		u := createUser("owner@example.com")
		desc := "My favorites"
		p := &playlist.Playlist{ID: uuid.New(), UserID: u.ID, Name: "Best Hits", Description: &desc}

		created, err := repo.CreatePlaylist(ctx, p)
		require.NoError(t, err)
		assert.Equal(t, p.Name, created.Name)

		found, err := repo.GetPlaylistByIDForUser(ctx, p.ID, u.ID)
		require.NoError(t, err)
		assert.Equal(t, p.ID, found.ID)

		t.Run("WrongUser", func(t *testing.T) {
			other := createUser("other@example.com")
			_, err := repo.GetPlaylistByIDForUser(ctx, p.ID, other.ID)
			assert.ErrorIs(t, err, playlist.ErrPlaylistNotFound)
		})
	})

	t.Run("TrackManagement", func(t *testing.T) {
		u := createUser("trackman@example.com")
		tr := createTrack("Song A")
		p, _ := repo.CreatePlaylist(ctx, &playlist.Playlist{ID: uuid.New(), UserID: u.ID, Name: "P1"})

		err := repo.AddTrackToPlaylist(ctx, p.ID, tr.ID, u.ID)
		require.NoError(t, err)

		tracks, err := repo.ListPlaylistTracks(ctx, p.ID, u.ID)
		require.NoError(t, err)
		assert.Len(t, tracks, 1)
		assert.Equal(t, tr.ID, tracks[0].ID)

		err = repo.RemoveTrackFromPlaylist(ctx, p.ID, tr.ID, u.ID)
		require.NoError(t, err)

		tracks, _ = repo.ListPlaylistTracks(ctx, p.ID, u.ID)
		assert.Empty(t, tracks)
	})

	t.Run("Update", func(t *testing.T) {
		u := createUser("updatep@example.com")
		p, _ := repo.CreatePlaylist(ctx, &playlist.Playlist{ID: uuid.New(), UserID: u.ID, Name: "Old Name"})

		p.Name = "New Name"
		updated, err := repo.UpdatePlaylist(ctx, p)
		require.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
	})

	t.Run("CountAndList", func(t *testing.T) {
		u := createUser("counter@example.com")
		repo.CreatePlaylist(ctx, &playlist.Playlist{ID: uuid.New(), UserID: u.ID, Name: "P1"})
		repo.CreatePlaylist(ctx, &playlist.Playlist{ID: uuid.New(), UserID: u.ID, Name: "P2"})

		count, err := repo.CountPlaylistsByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		list, err := repo.ListPlaylistsByUserID(ctx, u.ID)
		require.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("SecurityEdgeCases", func(t *testing.T) {
		u1 := createUser("u1@p.com")
		u2 := createUser("u2@p.com")
		p1, _ := repo.CreatePlaylist(ctx, &playlist.Playlist{ID: uuid.New(), UserID: u1.ID, Name: "U1 Playlist"})
		tr := createTrack("Song")

		t.Run("UpdateOtherUserPlaylist", func(t *testing.T) {
			p1.UserID = u2.ID // Attempting to change ownership or update as other user
			_, err := repo.UpdatePlaylist(ctx, p1)
			assert.ErrorIs(t, err, playlist.ErrPlaylistNotFound)
		})

		t.Run("AddTrackToNonExistentPlaylist", func(t *testing.T) {
			err := repo.AddTrackToPlaylist(ctx, uuid.New(), tr.ID, u1.ID)
			assert.Error(t, err)
		})

		t.Run("AddNonExistentTrackToPlaylist", func(t *testing.T) {
			err := repo.AddTrackToPlaylist(ctx, p1.ID, uuid.New(), u1.ID)
			assert.Error(t, err)
		})
	})
}
