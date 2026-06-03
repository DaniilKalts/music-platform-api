package favoriterepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/repository/testutil"
	trackrepo "github.com/DaniilKalts/music-platform-api/internal/repository/track"
	userrepo "github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

func TestFavoriteRepository(t *testing.T) {
	pool, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	q := sqlc.New(pool)
	repo := favoriterepo.NewRepository(q)
	uRepo := userrepo.NewRepository(pool)
	tRepo := trackrepo.NewRepository(q)
	ctx := context.Background()

	// Helpers
	u, _ := user.NewUser("fav@example.com", "fav")
	pass, _ := user.NewPassword("pass")
	createdU, _ := uRepo.Create(ctx, *u, pass)

	g, _ := tRepo.CreateGenre(ctx, &track.Genre{ID: uuid.New(), Name: "G"})
	ar, _ := tRepo.CreateArtist(ctx, &track.Artist{ID: uuid.New(), Name: "Ar"})
	al, _ := tRepo.CreateAlbum(ctx, &track.Album{ID: uuid.New(), Name: "Al"})
	tr, _ := tRepo.CreateTrack(ctx, &track.Track{ID: uuid.New(), Title: "T", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID})

	t.Run("AddAndCheck", func(t *testing.T) {
		err := repo.AddFavorite(ctx, &favorite.Favorite{UserID: createdU.ID, TrackID: tr.ID})
		require.NoError(t, err)

		exists, err := repo.FavoriteExists(ctx, createdU.ID, tr.ID)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("List", func(t *testing.T) {
		favs, err := repo.ListFavoritesByUserID(ctx, createdU.ID)
		require.NoError(t, err)
		assert.Len(t, favs, 1)
		assert.Equal(t, tr.ID, favs[0].ID)

		count, _ := repo.CountFavoritesByUserID(ctx, createdU.ID)
		assert.Equal(t, int64(1), count)
	})

	t.Run("NonExistent", func(t *testing.T) {
		exists, err := repo.FavoriteExists(ctx, createdU.ID, uuid.New())
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Remove", func(t *testing.T) {
		err := repo.RemoveFavorite(ctx, createdU.ID, tr.ID)
		require.NoError(t, err)

		exists, _ := repo.FavoriteExists(ctx, createdU.ID, tr.ID)
		assert.False(t, exists)
	})
}
