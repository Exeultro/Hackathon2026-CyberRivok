package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type APIError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return e.Code + ": " + e.Message
}

var (
	ErrUnauthorized  = &APIError{Code: "unauthorized", Message: "unauthorized"}
	ErrInvalidCreds  = &APIError{Code: "invalid_credentials", Message: "invalid email or password"}
	ErrEmailTaken    = &APIError{Code: "email_taken", Message: "email already taken"}
	ErrNicknameTaken = &APIError{Code: "nickname_taken", Message: "nickname already taken"}
	ErrModeLocked    = &APIError{Code: "mode_locked", Message: "mode or difficulty is locked"}
	ErrNotFound      = &APIError{Code: "not_found", Message: "not found"}
	ErrValidation    = &APIError{Code: "validation_error", Message: "validation error"}
	ErrRateLimited   = &APIError{Code: "rate_limited", Message: "too many requests"}
	ErrInternal      = &APIError{Code: "internal_error", Message: "internal server error"}
)

func StatusForCode(code string) int {
	switch code {
	case "unauthorized", "invalid_credentials":
		return http.StatusUnauthorized
	case "email_taken", "nickname_taken":
		return http.StatusConflict
	case "mode_locked":
		return http.StatusForbidden
	case "not_found":
		return http.StatusNotFound
	case "validation_error":
		return http.StatusUnprocessableEntity
	case "rate_limited":
		return http.StatusTooManyRequests
	case "not_implemented":
		return http.StatusNotImplemented
	case "method_not_allowed":
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		slog.ErrorContext(r.Context(), "unhandled error", "err", err)
		apiErr = ErrInternal
	}
	WriteJSON(w, StatusForCode(apiErr.Code), apiErr)
}

func WriteValidation(w http.ResponseWriter, r *http.Request, details map[string]any) {
	WriteError(w, r, &APIError{Code: "validation_error", Message: "validation error", Details: details})
}

type ctxKey string

const CtxUserID ctxKey = "user_id"

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxUserID, id)
}

func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(CtxUserID).(string)
	return id, ok && id != ""
}
