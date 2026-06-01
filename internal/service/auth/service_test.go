package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	"github.com/DaniilKalts/music-platform-api/pkg/jwt"
)

func TestServiceRegister(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		users := &userRepositoryMock{}
		service := NewService(users, &tokenManagerMock{})

		created, err := service.Register(context.Background(), RegisterInput{
			Email:    "Daniil.Kalts@Rbk.kz",
			Username: "daniilkalts",
			Password: "12345678",
		})

		require.NoError(t, err)
		require.Equal(t, "daniil.kalts@rbk.kz", created.Email)
		require.Equal(t, "daniilkalts", created.Username)
		require.Equal(t, user.RoleUser, created.Role)
		require.Equal(t, user.SubscriptionFree, created.Subscription)
		require.NotEmpty(t, users.createdPassword.Hash)
		require.NotEmpty(t, users.createdPassword.Salt)
	})

	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		service := NewService(&userRepositoryMock{}, &tokenManagerMock{})

		_, err := service.Register(context.Background(), RegisterInput{
			Email:    "daniil.kalts@rbk.kz",
			Username: "daniilkalts",
			Password: "short",
		})

		require.ErrorIs(t, err, user.ErrInvalidPassword)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		repositoryErr := errors.New("repository error")
		service := NewService(&userRepositoryMock{createErr: repositoryErr}, &tokenManagerMock{})

		_, err := service.Register(context.Background(), RegisterInput{
			Email:    "daniil.kalts@rbk.kz",
			Username: "daniilkalts",
			Password: "12345678",
		})

		require.ErrorIs(t, err, repositoryErr)
	})
}

func TestServiceLogin(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		password, err := user.NewPassword("12345678")
		require.NoError(t, err)

		userID := uuid.New()
		refreshExpiresAt := time.Now().Add(time.Hour)
		service := NewService(
			&userRepositoryMock{credentialsUser: &user.User{ID: userID, Email: "daniil.kalts@rbk.kz", Role: user.RoleUser}, credentialsPassword: password},
			&tokenManagerMock{pair: &jwt.Pair{
				AccessToken:           "access-token",
				AccessTokenExpiresAt:  time.Now().Add(time.Minute),
				RefreshToken:          "refresh-token",
				RefreshTokenExpiresAt: refreshExpiresAt,
			}},
		)

		pair, err := service.Login(context.Background(), LoginInput{Email: "DANIIL.KALTS@rbk.kz", Password: "12345678"})

		require.NoError(t, err)
		require.Equal(t, "access-token", pair.AccessToken)
		require.Equal(t, "refresh-token", pair.RefreshToken)
		require.Equal(t, refreshExpiresAt, pair.RefreshTokenExpiresAt)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		service := NewService(&userRepositoryMock{credentialsErr: user.ErrNotFound}, &tokenManagerMock{})

		_, err := service.Login(context.Background(), LoginInput{Email: "daniil.kalts@rbk.kz", Password: "12345678"})

		require.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()

		password, err := user.NewPassword("12345678")
		require.NoError(t, err)

		service := NewService(
			&userRepositoryMock{credentialsUser: &user.User{ID: uuid.New(), Email: "daniil.kalts@rbk.kz", Role: user.RoleUser}, credentialsPassword: password},
			&tokenManagerMock{},
		)

		_, err = service.Login(context.Background(), LoginInput{Email: "daniil.kalts@rbk.kz", Password: "wrong-password"})

		require.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

type userRepositoryMock struct {
	createErr error

	createdPassword user.Password

	credentialsUser     *user.User
	credentialsPassword user.Password
	credentialsErr      error
}

func (m *userRepositoryMock) Create(_ context.Context, u user.User, password user.Password) (*user.User, error) {
	m.createdPassword = password
	if m.createErr != nil {
		return nil, m.createErr
	}

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	return &u, nil
}

func (m *userRepositoryMock) GetCredentialsByEmail(_ context.Context, _ string) (*user.User, user.Password, error) {
	if m.credentialsErr != nil {
		return nil, user.Password{}, m.credentialsErr
	}

	return m.credentialsUser, m.credentialsPassword, nil
}

type tokenManagerMock struct {
	pair *jwt.Pair
	err  error
}

func (m *tokenManagerMock) GeneratePair(_ uuid.UUID, _ string) (*jwt.Pair, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.pair, nil
}
