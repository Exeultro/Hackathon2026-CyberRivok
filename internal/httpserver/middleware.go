package httpserver

import (
	"net/http"
	"strings"

	"typing/internal/httpx"
	"typing/internal/token"
)

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := s.parseBearer(r)
		if err != nil {
			httpx.WriteError(w, r, httpx.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), uid)))
	})
}

func (s *Server) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uid, err := s.parseBearer(r); err == nil {
			r = r.WithContext(httpx.WithUserID(r.Context(), uid))
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) parseBearer(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return "", httpx.ErrUnauthorized
	}
	claims, err := token.Parse(s.cfg.JWT.Secret, s.cfg.JWT.Issuer, strings.TrimPrefix(h, "Bearer "), token.TypeAccess)
	if err != nil {
		return "", httpx.ErrUnauthorized
	}
	return claims.Subject, nil
}
