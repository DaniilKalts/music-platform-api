package httpx

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/DaniilKalts/music-platform-api/pkg/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}

func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
	logger.FromContext(r.Context()).Error("unhandled service error", zap.Error(err))
	WriteError(w, http.StatusInternalServerError, "internal server error")
}
