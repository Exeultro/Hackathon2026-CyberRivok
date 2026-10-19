package maps

import (
	"context"
	"net/http"

	"typing/internal/httpx"
)

type ProgressProvider interface {
	Level(ctx context.Context, userID string) (int, error)
}

type Handler struct {
	svc      *Service
	progress ProgressProvider
}

func NewHandler(svc *Service, progress ProgressProvider) *Handler {
	return &Handler{svc: svc, progress: progress}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	uid, _ := httpx.UserID(r.Context())
	level := 1
	if uid != "" {
		l, err := h.progress.Level(r.Context(), uid)
		if err == nil && l > 0 {
			level = l
		}
	}
	items, err := h.svc.List(r.Context(), uid, level)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}
