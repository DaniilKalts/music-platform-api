package userrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/internal/repository/testutil"
	"github.com/DaniilKalts/music-platform-api/internal/repository/user"
)

func TestUserRepository(t *testing.T) {
	pool, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userrepo.NewRepository(pool)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		u, _ := user.NewUser("test@example.com", "testuser")
		pass, _ := user.NewPassword("password123")

		created, err := repo.Create(ctx, *u, pass)
		require.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, u.Email, created.Email)
		assert.Equal(t, u.Username, created.Username)

		t.Run("DuplicateEmail", func(t *testing.T) {
			u2, _ := user.NewUser("test@example.com", "otheruser")
			_, err := repo.Create(ctx, *u2, pass)
			assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
		})

		t.Run("DuplicateUsername", func(t *testing.T) {
			u2, _ := user.NewUser("other@example.com", "testuser")
			_, err := repo.Create(ctx, *u2, pass)
			assert.ErrorIs(t, err, user.ErrUsernameAlreadyExists)
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		u, _ := user.NewUser("getbyid@example.com", "getbyid")
		pass, _ := user.NewPassword("password123")
		created, _ := repo.Create(ctx, *u, pass)

		found, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.Email, found.Email)

		t.Run("NotFound", func(t *testing.T) {
			found, err := repo.GetByID(ctx, uuid.New())
			assert.ErrorIs(t, err, user.ErrNotFound)
			assert.Nil(t, found)
		})
	})

	t.Run("GetCredentialsByEmail", func(t *testing.T) {
		email := "creds@example.com"
		u, _ := user.NewUser(email, "credsuser")
		pass, _ := user.NewPassword("password123")
		repo.Create(ctx, *u, pass)

		foundU, foundPass, err := repo.GetCredentialsByEmail(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, u.Email, foundU.Email)
		assert.Equal(t, pass.Hash, foundPass.Hash)

		t.Run("NotFound", func(t *testing.T) {
			_, _, err := repo.GetCredentialsByEmail(ctx, "nonexistent@example.com")
			assert.ErrorIs(t, err, user.ErrNotFound)
		})
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		u1, _ := user.NewUser("u1@example.com", "u1")
		u2, _ := user.NewUser("u2@example.com", "u2")
		pass, _ := user.NewPassword("pass")
		created1, _ := repo.Create(ctx, *u1, pass)
		repo.Create(ctx, *u2, pass)

		t.Run("Success", func(t *testing.T) {
			newEmail := "updated@example.com"
			updated, err := repo.UpdateProfile(ctx, created1.ID, &newEmail, nil)
			require.NoError(t, err)
			assert.Equal(t, newEmail, updated.Email)
		})

		t.Run("DuplicateEmail", func(t *testing.T) {
			email2 := "u2@example.com"
			_, err := repo.UpdateProfile(ctx, created1.ID, &email2, nil)
			assert.ErrorIs(t, err, user.ErrEmailAlreadyExists)
		})

		t.Run("NotFound", func(t *testing.T) {
			email := "random@example.com"
			_, err := repo.UpdateProfile(ctx, uuid.New(), &email, nil)
			assert.ErrorIs(t, err, user.ErrNotFound)
		})
	})
}
