package history

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	domainhistory "github.com/DaniilKalts/music-platform-api/internal/domain/history"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	ListHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*domainhistory.HistoryRecord, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type HistoryRecordResponse struct {
	TrackID    string    `json:"track_id"`
	Title      string    `json:"title"`
	ArtistName string    `json:"artist_name"`
	ListenedAt time.Time `json:"listened_at"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt32(r, "limit", 20)
	offset := httpx.QueryInt32(r, "offset", 0)

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	records, err := h.service.ListHistory(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	res := make([]HistoryRecordResponse, len(records))
	for i, rec := range records {
		res[i] = HistoryRecordResponse{
			TrackID:    rec.TrackID.String(),
			Title:      rec.TrackTitle,
			ArtistName: rec.ArtistName,
			ListenedAt: rec.ListenedAt,
		}
	}

	httpx.JSON(w, http.StatusOK, res)
}
