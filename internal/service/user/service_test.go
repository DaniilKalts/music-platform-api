package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
)

func TestServiceGetMe(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	expected := &domainuser.User{ID: userID, Email: "daniil.kalts@rbk.kz", Username: "daniilkalts"}
	service := NewService(&repositoryMock{user: expected})

	actual, err := service.GetMe(context.Background(), userID)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestServiceUpdateMe(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	email := "New.Email@Rbk.kz"
	username := "newusername"
	repository := &repositoryMock{
		user: &domainuser.User{
			ID:           userID,
			Email:        "daniil.kalts@rbk.kz",
			Username:     "daniilkalts",
			Role:         domainuser.RoleUser,
			Subscription: domainuser.SubscriptionFree,
		},
	}
	service := NewService(repository)

	updated, err := service.UpdateMe(context.Background(), userID, UpdateInput{Email: &email, Username: &username})

	require.NoError(t, err)
	require.Equal(t, "new.email@rbk.kz", updated.Email)
	require.Equal(t, "newusername", updated.Username)
	require.Equal(t, "new.email@rbk.kz", *repository.updatedEmail)
	require.Equal(t, "newusername", *repository.updatedUsername)
}

func TestServiceUpdateMeValidationError(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	email := "bad"
	service := NewService(&repositoryMock{
		user: &domainuser.User{
			ID:           userID,
			Email:        "daniil.kalts@rbk.kz",
			Username:     "daniilkalts",
			Role:         domainuser.RoleUser,
			Subscription: domainuser.SubscriptionFree,
		},
	})

	_, err := service.UpdateMe(context.Background(), userID, UpdateInput{Email: &email})

	require.ErrorIs(t, err, domainuser.ErrInvalidEmail)
}

func TestServiceUpdateMeRepositoryError(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("repository error")
	service := NewService(&repositoryMock{getErr: repositoryErr})

	_, err := service.UpdateMe(context.Background(), uuid.New(), UpdateInput{})

	require.ErrorIs(t, err, repositoryErr)
}

type repositoryMock struct {
	user *domainuser.User

	getErr    error
	updateErr error

	updatedEmail    *string
	updatedUsername *string
}

func (m *repositoryMock) GetByID(_ context.Context, _ uuid.UUID) (*domainuser.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.user, nil
}

func (m *repositoryMock) UpdateProfile(_ context.Context, _ uuid.UUID, email, username *string) (*domainuser.User, error) {
	m.updatedEmail = email
	m.updatedUsername = username
	if m.updateErr != nil {
		return nil, m.updateErr
	}

	updated := *m.user
	if email != nil {
		updated.Email = domainuser.NormalizeEmail(*email)
	}
	if username != nil {
		updated.Username = *username
	}
	return &updated, nil
}
