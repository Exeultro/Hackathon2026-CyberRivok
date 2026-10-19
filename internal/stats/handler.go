package stats

import (
	"net/http"
	"strconv"
	"time"

	"typing/internal/httpx"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	s, err := h.repo.Summary(r.Context(), uid)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, s)
}

func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	rangeParam := r.URL.Query().Get("range")
	var from time.Time
	switch rangeParam {
	case "week", "":
		rangeParam = "week"
		from = time.Now().UTC().AddDate(0, 0, -7)
	case "month":
		from = time.Now().UTC().AddDate(0, -1, 0)
	case "year":
		from = time.Now().UTC().AddDate(-1, 0, 0)
	default:
		httpx.WriteValidation(w, r, map[string]any{"range": "must be one of: week, month, year"})
		return
	}

	points, err := h.repo.Timeline(r.Context(), uid, from)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"range": rangeParam, "points": points})
}

func (h *Handler) Keys(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	kind := r.URL.Query().Get("kind")
	if kind == "" {
		kind = "key"
	}
	if kind != "key" && kind != "bigram" {
		httpx.WriteValidation(w, r, map[string]any{"kind": "must be one of: key, bigram"})
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang != "" && lang != "ru" && lang != "en" {
		httpx.WriteValidation(w, r, map[string]any{"lang": "must be one of: ru, en"})
		return
	}
	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	items, err := h.repo.KeyStats(r.Context(), uid, kind, lang, limit)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"kind": kind, "items": items})
}
