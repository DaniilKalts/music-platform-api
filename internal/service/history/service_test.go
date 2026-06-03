package history_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
	service "github.com/DaniilKalts/music-platform-api/internal/service/history"
)

type mockHistoryRepo struct{ mock.Mock }

func (m *mockHistoryRepo) ListListeningHistoryByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*history.HistoryRecord, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*history.HistoryRecord), args.Error(1)
}

func TestListHistory(t *testing.T) {
	ctx := context.Background()
	uID := uuid.New()
	limit, offset := int32(10), int32(0)

	t.Run("Success", func(t *testing.T) {
		mRepo := new(mockHistoryRepo)
		s := service.NewService(mRepo)

		expected := []*history.HistoryRecord{
			{ID: uuid.New(), UserID: uID, TrackID: uuid.New(), ListenedAt: time.Now()},
		}
		mRepo.On("ListListeningHistoryByUserID", ctx, uID, limit, offset).Return(expected, nil)

		res, err := s.ListHistory(ctx, uID, limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
