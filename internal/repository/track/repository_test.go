package trackrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/internal/repository/testutil"
	"github.com/DaniilKalts/music-platform-api/internal/repository/track"
)

func TestTrackRepository(t *testing.T) {
	pool, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := trackrepo.NewRepository(pool)
	ctx := context.Background()

	createGenre := func(name string) *track.Genre {
		g := &track.Genre{ID: uuid.New(), Name: name}
		created, err := repo.CreateGenre(ctx, g)
		require.NoError(t, err)
		return created
	}

	createArtist := func(name string) *track.Artist {
		a := &track.Artist{ID: uuid.New(), Name: name}
		created, err := repo.CreateArtist(ctx, a)
		require.NoError(t, err)
		return created
	}

	createAlbum := func(name string) *track.Album {
		al := &track.Album{ID: uuid.New(), Name: name}
		created, err := repo.CreateAlbum(ctx, al)
		require.NoError(t, err)
		return created
	}

	t.Run("CreateTrack", func(t *testing.T) {
		g := createGenre("Rock")
		ar := createArtist("Linkin Park")
		al := createAlbum("Meteora")

		tr := &track.Track{
			ID:              uuid.New(),
			Title:           "Numb",
			ArtistID:        ar.ID,
			AlbumID:         al.ID,
			GenreID:         g.ID,
			DurationSeconds: 187,
			FileURL:         "https://example.com/numb.mp3",
		}

		created, err := repo.CreateTrack(ctx, tr)
		require.NoError(t, err)
		assert.Equal(t, tr.Title, created.Title)
		assert.Equal(t, ar.Name, created.ArtistName)
		assert.Equal(t, al.Name, created.AlbumName)
	})

	t.Run("GetTrackByID", func(t *testing.T) {
		g := createGenre("Pop")
		ar := createArtist("Taylor Swift")
		al := createAlbum("1989")
		tr := &track.Track{ID: uuid.New(), Title: "Style", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID, DurationSeconds: 231, FileURL: "url"}
		repo.CreateTrack(ctx, tr)

		found, err := repo.GetTrackByID(ctx, tr.ID)
		require.NoError(t, err)
		assert.Equal(t, tr.Title, found.Title)

		t.Run("NotFound", func(t *testing.T) {
			_, err := repo.GetTrackByID(ctx, uuid.New())
			assert.ErrorIs(t, err, track.ErrTrackNotFound)
		})
	})

	t.Run("ListTracks", func(t *testing.T) {
		g := createGenre("Jazz")
		ar := createArtist("Miles Davis")
		al := createAlbum("Kind of Blue")

		for i := 0; i < 5; i++ {
			tr := &track.Track{ID: uuid.New(), Title: "Track", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID}
			repo.CreateTrack(ctx, tr)
		}

		tracks, err := repo.ListTracks(ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, tracks, 3)
	})

	t.Run("SearchTracks", func(t *testing.T) {
		g := createGenre("Electronic")
		ar := createArtist("Daft Punk")
		al := createAlbum("Discovery")
		tr := &track.Track{ID: uuid.New(), Title: "One More Time", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID}
		repo.CreateTrack(ctx, tr)

		tracks, err := repo.SearchTracks(ctx, "Daft", 10, 0)
		require.NoError(t, err)
		assert.NotEmpty(t, tracks)
		assert.Equal(t, "One More Time", tracks[0].Title)
	})

	t.Run("SoftDelete", func(t *testing.T) {
		g := createGenre("Metal")
		ar := createArtist("Metallica")
		al := createAlbum("Master of Puppets")
		tr := &track.Track{ID: uuid.New(), Title: "Battery", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID}
		repo.CreateTrack(ctx, tr)

		err := repo.SoftDeleteTrack(ctx, tr.ID)
		require.NoError(t, err)

		_, err = repo.GetTrackByID(ctx, tr.ID)
		assert.ErrorIs(t, err, track.ErrTrackNotFound)

		exists, _ := repo.TrackExists(ctx, tr.ID)
		assert.False(t, exists)
	})

	t.Run("UpdateTrack", func(t *testing.T) {
		g := createGenre("Electronic")
		ar := createArtist("Daft Punk")
		al := createAlbum("Homework")
		tr := &track.Track{ID: uuid.New(), Title: "Around the World", ArtistID: ar.ID, AlbumID: al.ID, GenreID: g.ID}
		repo.CreateTrack(ctx, tr)

		tr.Title = "Around the World (Edit)"
		updated, err := repo.UpdateTrack(ctx, tr)
		require.NoError(t, err)
		assert.Equal(t, "Around the World (Edit)", updated.Title)

		t.Run("NotFound", func(t *testing.T) {
			tr2 := &track.Track{ID: uuid.New(), Title: "X"}
			_, err := repo.UpdateTrack(ctx, tr2)
			assert.ErrorIs(t, err, track.ErrTrackNotFound)
		})
	})

	t.Run("ListGenres", func(t *testing.T) {
		createGenre("Genre A")
		createGenre("Genre B")
		genres, err := repo.ListGenres(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(genres), 2)
	})

	t.Run("AlbumIdempotency", func(t *testing.T) {
		al := &track.Album{ID: uuid.New(), Name: "Discovery"}
		c1, err := repo.FindOrCreateAlbum(ctx, al)
		require.NoError(t, err)

		c2, err := repo.FindOrCreateAlbum(ctx, al)
		require.NoError(t, err)
		assert.Equal(t, c1.ID, c2.ID)
	})
}
