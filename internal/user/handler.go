package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"typing/internal/httpx"
	"typing/internal/settings"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type userDTO struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Nickname  string            `json:"nickname"`
	CreatedAt string            `json:"created_at"`
	Settings  settings.Settings `json:"settings"`
}

type updateReq struct {
	Nickname *string         `json:"nickname,omitempty"`
	Settings *settings.Patch `json:"settings,omitempty"`
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}
	u, err := h.svc.GetMe(r.Context(), uid)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, r, httpx.ErrNotFound)
			return
		}
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toDTO(u))
}

func (h *Handler) PatchMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := httpx.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, r, httpx.ErrUnauthorized)
		return
	}

	var req updateReq
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httpx.WriteValidation(w, r, map[string]any{"body": "invalid json: " + err.Error()})
		return
	}

	in := UpdateInput{Nickname: req.Nickname, Settings: req.Settings}
	if details := ValidateUpdate(in); len(details) > 0 {
		httpx.WriteValidation(w, r, details)
		return
	}

	u, err := h.svc.UpdateMe(r.Context(), uid, in)
	if err != nil {
		switch {
		case errors.Is(err, ErrNicknameTaken):
			httpx.WriteError(w, r, httpx.ErrNicknameTaken)
		case errors.Is(err, ErrNotFound):
			httpx.WriteError(w, r, httpx.ErrNotFound)
		default:
			httpx.WriteError(w, r, err)
		}
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toDTO(u))
}

func toDTO(u *User) userDTO {
	return userDTO{
		ID:        u.ID,
		Email:     u.Email,
		Nickname:  u.Nickname,
		CreatedAt: u.CreatedAt,
		Settings:  u.Settings,
	}
}
