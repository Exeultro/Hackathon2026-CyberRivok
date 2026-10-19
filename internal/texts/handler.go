package texts

import (
	"errors"
	"net/http"
	"strconv"

	"typing/internal/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	mode := q.Get("mode")
	if mode == "" {
		httpx.WriteValidation(w, r, map[string]any{"mode": "required"})
		return
	}
	lang := q.Get("lang")
	if lang == "" {
		httpx.WriteValidation(w, r, map[string]any{"lang": "required"})
		return
	}

	difficulty := 1
	if v := q.Get("difficulty"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 5 {
			httpx.WriteValidation(w, r, map[string]any{"difficulty": "must be 1..5"})
			return
		}
		difficulty = n
	}

	count := 1
	if v := q.Get("count"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 50 {
			httpx.WriteValidation(w, r, map[string]any{"count": "must be 1..50"})
			return
		}
		count = n
	}

	uid, _ := httpx.UserID(r.Context())

	items, err := h.svc.Get(r.Context(), Request{
		Mode:         mode,
		Lang:         lang,
		Difficulty:   difficulty,
		Count:        count,
		CodeLanguage: q.Get("code_language"),
		UserID:       uid,
	})
	if err != nil {
		if errors.Is(err, ErrValidation) {
			httpx.WriteValidation(w, r, map[string]any{"query": "invalid parameters"})
			return
		}
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}
