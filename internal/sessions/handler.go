package sessions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"typing/internal/httpx"
	"typing/internal/progress"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}

	var req CreateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httpx.WriteValidation(w, r, map[string]any{"body": "invalid json: " + err.Error()})
		return
	}

	res, err := h.svc.Create(r.Context(), uid, req)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			httpx.WriteValidation(w, r, map[string]any{"body": "invalid session payload"})
			return
		}
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toSessionResult(res))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	limit := parseIntDefault(r.URL.Query().Get("limit"), 20, 1, 100)
	offset := parseIntDefault(r.URL.Query().Get("offset"), 0, 0, 1<<30)

	items, total, err := h.svc.List(r.Context(), uid, limit, offset)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	out := map[string]any{
		"items": sessionSummaries(items),
		"total": total,
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

func toSessionResult(res *Result) map[string]any {
	s := res.Session
	out := map[string]any{
		"id":               s.ID,
		"mode":             s.Mode,
		"lang":             s.Lang,
		"difficulty":       s.Difficulty,
		"started_at":       s.StartedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"duration_ms":      s.DurationMS,
		"wpm":              s.WPM,
		"cpm":              s.CPM,
		"accuracy":         s.Accuracy,
		"errors":           s.Errors,
		"score":            s.Score,
		"xp_gained":        s.XPGained,
		"level_before":     res.LevelBefore,
		"progress":         res.Progress,
		"new_achievements": res.NewAchievements,
		"unlocked":         res.Unlocked,
		"is_personal_best": res.IsPersonalBest,
	}
	if res.NodeResult != nil {
		out["node_result"] = res.NodeResult
	}
	if res.NewAchievements == nil {
		out["new_achievements"] = []progress.Achievement{}
	}
	if res.Unlocked == nil {
		out["unlocked"] = []progress.Unlock{}
	}
	return out
}

func sessionSummaries(items []Session) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, s := range items {
		out = append(out, map[string]any{
			"id":          s.ID,
			"mode":        s.Mode,
			"lang":        s.Lang,
			"difficulty":  s.Difficulty,
			"started_at":  s.StartedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			"duration_ms": s.DurationMS,
			"wpm":         s.WPM,
			"cpm":         s.CPM,
			"accuracy":    s.Accuracy,
			"errors":      s.Errors,
			"score":       s.Score,
			"xp_gained":   s.XPGained,
		})
	}
	return out
}

func parseIntDefault(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < min || n > max {
		return def
	}
	return n
}
