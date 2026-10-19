package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"typing/internal/httpx"
	"typing/internal/settings"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type tokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type userDTO struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Nickname  string            `json:"nickname"`
	CreatedAt string            `json:"created_at"`
	Settings  settings.Settings `json:"settings"`
}

type authResponse struct {
	Tokens tokensDTO `json:"tokens"`
	User   userDTO   `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if !decodeJSON(w, r, &req) {
		return
	}

	details := map[string]any{}
	if !validEmail(req.Email) {
		details["email"] = "invalid email"
	}
	if l := utf8.RuneCountInString(req.Password); l < 8 || l > 128 {
		details["password"] = "length must be 8..128"
	}
	if l := utf8.RuneCountInString(req.Nickname); l < 2 || l > 24 {
		details["nickname"] = "length must be 2..24"
	}
	if len(details) > 0 {
		httpx.WriteValidation(w, r, details)
		return
	}

	res, err := h.svc.Register(r.Context(), req.Email, req.Password, req.Nickname)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			httpx.WriteError(w, r, httpx.ErrEmailTaken)
		case errors.Is(err, ErrNicknameTaken):
			httpx.WriteError(w, r, httpx.ErrNicknameTaken)
		default:
			httpx.WriteError(w, r, err)
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toAuthResponse(res))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Email == "" || req.Password == "" {
		httpx.WriteValidation(w, r, map[string]any{"email": "required", "password": "required"})
		return
	}

	res, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpx.WriteError(w, r, httpx.ErrInvalidCreds)
			return
		}
		httpx.WriteError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toAuthResponse(res))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if !decodeJSON(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		httpx.WriteValidation(w, r, map[string]any{"refresh_token": "required"})
		return
	}

	res, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpx.WriteError(w, r, httpx.ErrUnauthorized)
			return
		}
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, tokensDTO{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    res.ExpiresIn,
	})
}

func toAuthResponse(res *AuthResult) authResponse {
	return authResponse{
		Tokens: tokensDTO{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			TokenType:    "bearer",
			ExpiresIn:    res.ExpiresIn,
		},
		User: userDTO{
			ID:        res.User.ID,
			Email:     res.User.Email,
			Nickname:  res.User.Nickname,
			CreatedAt: res.User.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			Settings:  settings.Default(),
		},
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		httpx.WriteValidation(w, r, map[string]any{"body": "invalid json: " + err.Error()})
		return false
	}
	return true
}

func validEmail(s string) bool {
	at := strings.LastIndexByte(s, '@')
	if at <= 0 || at == len(s)-1 {
		return false
	}
	domain := s[at+1:]
	return strings.Contains(domain, ".") && !strings.ContainsAny(s, " \t\r\n")
}
