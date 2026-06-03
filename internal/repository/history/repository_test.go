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

	// Helpers
	u, _ := user.NewUser("hist@example.com", "hist")
	pass, _ := user.NewPassword("pass")
	createdU, _ := uRepo.Create(ctx, *u, pass)

	g, _ := tRepo.CreateGenre(ctx, &track.Genre{ID: uuid.New(), Name: "G"})
	ar, _ := tRepo.CreateArtist(ctx, &track.Artist{ID: uuid.New(), Name: "Ar"})
	al, _ := tRepo.CreateAlbum(ctx, &track.Album{ID: uuid.New(), Name: "Al"})
	tr, _ := tRepo.CreateTrack(ctx, &track.Track{ID: uuid.New(), Title: "T", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID})

	t.Run("CreateAndList", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			h := &history.HistoryRecord{
				ID:      uuid.New(),
				UserID:  createdU.ID,
				TrackID: tr.ID,
			}
			repo.CreateListeningHistory(ctx, h)
		}

		list, err := repo.ListListeningHistoryByUserID(ctx, createdU.ID, 3, 0)
		require.NoError(t, err)
		assert.Len(t, list, 3)

		list2, err := repo.ListListeningHistoryByUserID(ctx, createdU.ID, 3, 3)
		require.NoError(t, err)
		assert.Len(t, list2, 2)
	})
}
