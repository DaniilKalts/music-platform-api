package admin

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	domaintrack "github.com/DaniilKalts/music-platform-api/internal/domain/track"
	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceadmin "github.com/DaniilKalts/music-platform-api/internal/service/admin"
	tracktransport "github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/track"
	usertransport "github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/user"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	CreateTrack(ctx context.Context, input serviceadmin.CreateTrackInput) (*domaintrack.Track, error)
	UpdateTrack(ctx context.Context, input serviceadmin.UpdateTrackInput) (*domaintrack.Track, error)
	DeleteTrack(ctx context.Context, id uuid.UUID) error
	UpdateUserSubscription(ctx context.Context, id uuid.UUID, sub domainuser.Subscription) (*domainuser.User, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	var req CreateTrackRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	genreID, err := uuid.Parse(req.GenreID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid genre id")
		return
	}

	t, err := h.service.CreateTrack(r.Context(), serviceadmin.CreateTrackInput{
		Title:           req.Title,
		ArtistName:      req.ArtistName,
		AlbumName:       req.AlbumName,
		GenreID:         genreID,
		DurationSeconds: req.DurationSeconds,
		FileURL:         req.FileURL,
	})
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, tracktransport.ToTrackResponse(*t))
}

func (h *Handler) UpdateTrack(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	var req UpdateTrackRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	genreID, err := uuid.Parse(req.GenreID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid genre id")
		return
	}

	t, err := h.service.UpdateTrack(r.Context(), serviceadmin.UpdateTrackInput{
		ID:              id,
		Title:           req.Title,
		ArtistName:      req.ArtistName,
		AlbumName:       req.AlbumName,
		GenreID:         genreID,
		DurationSeconds: req.DurationSeconds,
		FileURL:         req.FileURL,
	})
	if err != nil {
		if errors.Is(err, domaintrack.ErrTrackNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, tracktransport.ToTrackResponse(*t))
}

func (h *Handler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	if err := h.service.DeleteTrack(r.Context(), id); err != nil {
		if errors.Is(err, domaintrack.ErrTrackNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateUserSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateSubscriptionRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	u, err := h.service.UpdateUserSubscription(r.Context(), id, req.Type)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, usertransport.ToResponse(*u))
}
