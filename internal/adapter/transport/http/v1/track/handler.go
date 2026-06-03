package track

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	domaintrack "github.com/DaniilKalts/music-platform-api/internal/domain/track"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	GetTrack(ctx context.Context, id uuid.UUID) (*domaintrack.Track, error)
	ListTracks(ctx context.Context, limit, offset int32) ([]*domaintrack.Track, error)
	SearchTracks(ctx context.Context, query string, limit, offset int32) ([]*domaintrack.Track, error)
	ListGenres(ctx context.Context) ([]*domaintrack.Genre, error)
	PlayTrack(ctx context.Context, userID, trackID uuid.UUID) (*domaintrack.Track, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := httpx.QueryInt32(r, "limit", 20)
	offset := httpx.QueryInt32(r, "offset", 0)

	tracks, err := h.service.ListTracks(r.Context(), limit, offset)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToTracksResponse(tracks))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	t, err := h.service.GetTrack(r.Context(), id)
	if err != nil {
		if errors.Is(err, domaintrack.ErrTrackNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToTrackResponse(*t))
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit := httpx.QueryInt32(r, "limit", 20)
	offset := httpx.QueryInt32(r, "offset", 0)

	if query == "" {
		httpx.JSON(w, http.StatusOK, []TrackResponse{})
		return
	}

	tracks, err := h.service.SearchTracks(r.Context(), query, limit, offset)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToTracksResponse(tracks))
}

func (h *Handler) ListGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.service.ListGenres(r.Context())
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToGenresResponse(genres))
}

func (h *Handler) Play(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	t, err := h.service.PlayTrack(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domaintrack.ErrTrackNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToTrackResponse(*t))
}
