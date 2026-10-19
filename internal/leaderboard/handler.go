package leaderboard

import (
	"net/http"
	"strconv"

	"typing/internal/httpx"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	metric := q.Get("metric")
	if metric != "speed" && metric != "accuracy" && metric != "xp" {
		httpx.WriteValidation(w, r, map[string]any{"metric": "must be one of: speed, accuracy, xp"})
		return
	}
	period := q.Get("period")
	if period == "" {
		period = "all"
	}
	if period != "today" && period != "week" && period != "all" {
		httpx.WriteValidation(w, r, map[string]any{"period": "must be one of: today, week, all"})
		return
	}
	mode := q.Get("mode")
	if mode != "" && mode != "words" && mode != "sentences" && mode != "code" && mode != "symbols" {
		httpx.WriteValidation(w, r, map[string]any{"mode": "invalid"})
		return
	}
	lang := q.Get("lang")
	if lang != "" && lang != "ru" && lang != "en" {
		httpx.WriteValidation(w, r, map[string]any{"lang": "invalid"})
		return
	}
	limit := 20
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	offset := 0
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	f := Filter{Metric: metric, Period: period, Mode: mode, Lang: lang, Limit: limit, Offset: offset}
	items, total, err := h.repo.List(r.Context(), f)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	type itemDTO struct {
		Rank  int       `json:"rank"`
		User  EntryUser `json:"user"`
		Value float64   `json:"value"`
	}
	out := make([]itemDTO, 0, len(items))
	for _, e := range items {
		out = append(out, itemDTO{
			Rank:  e.Rank,
			User:  EntryUser{ID: e.UserID, Nickname: e.Nickname, Level: e.Level},
			Value: e.Value,
		})
	}

	resp := map[string]any{
		"metric": metric,
		"period": period,
		"items":  out,
		"total":  total,
	}

	if uid, ok := httpx.UserID(r.Context()); ok {
		me, err := h.repo.Me(r.Context(), f, uid)
		if err != nil {
			httpx.WriteError(w, r, err)
			return
		}
		if me != nil {
			resp["me"] = itemDTO{
				Rank:  me.Rank,
				User:  EntryUser{ID: me.UserID, Nickname: me.Nickname, Level: me.Level},
				Value: me.Value,
			}
		} else {
			resp["me"] = nil
		}
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
