package playlist

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	domainplaylist "github.com/DaniilKalts/music-platform-api/internal/domain/playlist"
	domaintrack "github.com/DaniilKalts/music-platform-api/internal/domain/track"
	serviceplaylist "github.com/DaniilKalts/music-platform-api/internal/service/playlist"
	tracktransport "github.com/DaniilKalts/music-platform-api/internal/adapter/transport/http/v1/track"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	CreatePlaylist(ctx context.Context, input serviceplaylist.CreateInput) (*domainplaylist.Playlist, error)
	GetPlaylist(ctx context.Context, id, userID uuid.UUID) (*domainplaylist.Playlist, error)
	ListPlaylists(ctx context.Context, userID uuid.UUID) ([]*domainplaylist.Playlist, error)
	UpdatePlaylist(ctx context.Context, input serviceplaylist.UpdateInput) (*domainplaylist.Playlist, error)
	DeletePlaylist(ctx context.Context, id, userID uuid.UUID) error
	AddTrack(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	RemoveTrack(ctx context.Context, playlistID, trackID, userID uuid.UUID) error
	ListTracks(ctx context.Context, playlistID, userID uuid.UUID) ([]*domaintrack.Track, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePlaylistRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	p, err := h.service.CreatePlaylist(r.Context(), serviceplaylist.CreateInput{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistLimitReached) {
			httpx.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, ToPlaylistResponse(*p))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	playlists, err := h.service.ListPlaylists(r.Context(), userID)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToPlaylistsResponse(playlists))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	p, err := h.service.GetPlaylist(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToPlaylistResponse(*p))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	var req UpdatePlaylistRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	p, err := h.service.UpdatePlaylist(r.Context(), serviceplaylist.UpdateInput{
		ID:          id,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToPlaylistResponse(*p))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.DeletePlaylist(r.Context(), id, userID); err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddTrack(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(chi.URLParam(r, "playlistID"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	trackID, err := uuid.Parse(chi.URLParam(r, "trackID"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.AddTrack(r.Context(), playlistID, trackID, userID); err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistNotFound) || errors.Is(err, domaintrack.ErrTrackNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(chi.URLParam(r, "playlistID"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	trackID, err := uuid.Parse(chi.URLParam(r, "trackID"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid track id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.RemoveTrack(r.Context(), playlistID, trackID, userID); err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTracks(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid playlist id")
		return
	}

	userID, err := httpx.ExtractUserID(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tracks, err := h.service.ListTracks(r.Context(), playlistID, userID)
	if err != nil {
		if errors.Is(err, domainplaylist.ErrPlaylistNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, tracktransport.ToTracksResponse(tracks))
}
