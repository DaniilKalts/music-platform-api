package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	domainuser "github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceuser "github.com/DaniilKalts/music-platform-api/internal/service/user"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	GetMe(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
	UpdateMe(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*domainuser.User, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	identity, ok := httpx.UserFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.service.GetMe(r.Context(), identity.ID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToResponse(*profile))
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	identity, ok := httpx.UserFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body UpdateRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	updated, err := h.service.UpdateMe(r.Context(), identity.ID, ToUpdateInput(body))
	if err != nil {
		writeError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, ToResponse(*updated))
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domainuser.ErrInvalidEmail), errors.Is(err, domainuser.ErrInvalidUsername):
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domainuser.ErrEmailAlreadyExists), errors.Is(err, domainuser.ErrUsernameAlreadyExists):
		httpx.WriteError(w, http.StatusConflict, err.Error())
	default:
		httpx.WriteInternalError(w, r, err)
	}
}
