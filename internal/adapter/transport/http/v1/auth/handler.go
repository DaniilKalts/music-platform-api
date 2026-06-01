package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/DaniilKalts/music-platform-api/internal/domain/user"
	serviceauth "github.com/DaniilKalts/music-platform-api/internal/service/auth"
	"github.com/DaniilKalts/music-platform-api/pkg/httpx"
)

type Service interface {
	Register(ctx context.Context, input serviceauth.RegisterInput) (*user.User, error)
	Login(ctx context.Context, input serviceauth.LoginInput) (*serviceauth.TokenPair, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	created, err := h.service.Register(r.Context(), ToRegisterInput(body))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidUsername),
			errors.Is(err, user.ErrInvalidPassword):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, user.ErrEmailAlreadyExists),
			errors.Is(err, user.ErrUsernameAlreadyExists):
			httpx.WriteError(w, http.StatusConflict, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, ToRegisterResponse(*created))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if !httpx.DecodeJSON(w, r, &body) {
		return
	}

	token, err := h.service.Login(r.Context(), ToLoginInput(body))
	if err != nil {
		switch {
		case errors.Is(err, serviceauth.ErrInvalidCredentials):
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteInternalError(w, r, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, ToTokenResponse(*token))
}
