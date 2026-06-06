package favorite

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	tracktransport "github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/track"
	"github.com/DaniilKalts/music-platform-api/internal/domain/favorite"
	"github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	AddFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, trackID uuid.UUID) error
	ListFavorites(ctx context.Context, userID uuid.UUID) ([]*track.Track, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	trackIDStr := chi.URLParam(r, "track_id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.AddFavorite(r.Context(), userID, trackID); err != nil {
		switch {
		case errors.Is(err, favorite.ErrFavoriteLimitReached):
			httpx.WriteError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, track.ErrTrackNotFound):
			httpx.WriteError(w, http.StatusNotFound, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	trackIDStr := chi.URLParam(r, "track_id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.RemoveFavorite(r.Context(), userID, trackID); err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tracks, err := h.service.ListFavorites(r.Context(), userID)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, tracktransport.ToTracksResponse(tracks))
}
