package historyrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository/history"
	"github.com/DaniilKalts/music-platform-api/internal/repository/testutil"
	trackrepo "github.com/DaniilKalts/music-platform-api/internal/repository/track"
	userrepo "github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

func TestHistoryRepository(t *testing.T) {
	pool, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := historyrepo.NewRepository(pool)
	uRepo := userrepo.NewRepository(pool)
	tRepo := trackrepo.NewRepository(pool)
	ctx := context.Background()

	u, err := user.NewUser("hist@example.com", "hist")
	require.NoError(t, err)
	pass, err := user.NewPassword("Password123!")
	require.NoError(t, err)
	createdU, err := uRepo.Create(ctx, *u, pass)
	require.NoError(t, err)

	g, err := tRepo.CreateGenre(ctx, &track.Genre{ID: uuid.New(), Name: "G"})
	require.NoError(t, err)
	ar, err := tRepo.CreateArtist(ctx, &track.Artist{ID: uuid.New(), Name: "Ar"})
	require.NoError(t, err)
	al, err := tRepo.CreateAlbum(ctx, &track.Album{ID: uuid.New(), Name: "Al"})
	require.NoError(t, err)
	tr, err := tRepo.CreateTrack(ctx, &track.Track{ID: uuid.New(), Title: "T", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID, DurationSeconds: 180})
	require.NoError(t, err)
	tr2, err := tRepo.CreateTrack(ctx, &track.Track{ID: uuid.New(), Title: "T2", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID, DurationSeconds: 180})
	require.NoError(t, err)

	t.Run("CreateAndList", func(t *testing.T) {
		for range 5 {
			h := &history.HistoryRecord{
				ID:      uuid.New(),
				UserID:  createdU.ID,
				TrackID: tr.ID,
			}
			repo.CreateListeningHistory(ctx, h)
		}
		repo.CreateListeningHistory(ctx, &history.HistoryRecord{
			ID:      uuid.New(),
			UserID:  createdU.ID,
			TrackID: tr2.ID,
		})

		list, err := repo.ListListeningHistoryByUserID(ctx, createdU.ID, 10, 0)
		require.NoError(t, err)
		require.Len(t, list, 2)
		assert.NotEqual(t, list[0].TrackID, list[1].TrackID)

		list2, err := repo.ListListeningHistoryByUserID(ctx, createdU.ID, 1, 1)
		require.NoError(t, err)
		assert.Len(t, list2, 1)

		list3, err := repo.ListListeningHistoryByUserID(ctx, createdU.ID, 10, 2)
		require.NoError(t, err)
		assert.Empty(t, list3)
	})
}
