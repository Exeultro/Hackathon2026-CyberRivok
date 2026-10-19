package progress

import (
	"net/http"

	"typing/internal/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProgress(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	p, err := h.svc.GetProgress(r.Context(), uid)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) ListAchievements(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		items, err := h.svc.ListAchievementsAnonymous(r.Context())
		if err != nil {
			httpx.WriteError(w, r, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	items, err := h.svc.ListAchievements(r.Context(), uid)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) DeleteProgress(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	if err := h.svc.Reset(r.Context(), uid); err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
