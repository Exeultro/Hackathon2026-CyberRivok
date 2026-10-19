package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"typing/internal/auth"
	"typing/internal/config"
	"typing/internal/httpx"
	"typing/internal/leaderboard"
	"typing/internal/maps"
	"typing/internal/progress"
	"typing/internal/sessions"
	"typing/internal/stats"
	"typing/internal/texts"
	"typing/internal/user"
)

type Server struct {
	cfg         *config.Config
	auth        *auth.Handler
	user        *user.Handler
	sessions    *sessions.Handler
	progress    *progress.Handler
	stats       *stats.Handler
	leaderboard *leaderboard.Handler
	texts       *texts.Handler
	maps        *maps.Handler
}

func New(
	cfg *config.Config,
	authH *auth.Handler,
	userH *user.Handler,
	sessH *sessions.Handler,
	progH *progress.Handler,
	statsH *stats.Handler,
	lbH *leaderboard.Handler,
	textsH *texts.Handler,
	mapsH *maps.Handler,
) *Server {
	return &Server{
		cfg:         cfg,
		auth:        authH,
		user:        userH,
		sessions:    sessH,
		progress:    progH,
		stats:       statsH,
		leaderboard: lbH,
		texts:       textsH,
		maps:        mapsH,
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(s.requestLogger())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Use(rateLimitAuth(10, time.Minute))
			r.Post("/register", s.auth.Register)
			r.Post("/login", s.auth.Login)
			r.Post("/refresh", s.auth.Refresh)
		})

		// public
		r.With(s.OptionalAuth).Get("/texts", s.texts.Get)
		r.With(s.OptionalAuth).Get("/achievements", s.progress.ListAchievements)
		r.With(s.OptionalAuth).Get("/leaderboard", s.leaderboard.Get)
		r.With(s.OptionalAuth).Get("/maps", s.maps.List)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(s.Auth)

			r.Get("/me", s.user.GetMe)
			r.Patch("/me", s.user.PatchMe)

			r.Post("/sessions", s.sessions.Create)
			r.Get("/sessions", s.sessions.List)
			r.Get("/me/progress", s.progress.GetProgress)
			r.Delete("/me/progress", s.progress.DeleteProgress)
			r.Get("/me/stats/summary", s.stats.Summary)
			r.Get("/me/stats/timeline", s.stats.Timeline)
			r.Get("/me/stats/keys", s.stats.Keys)
		})
	})

	if s.cfg.StaticDir != "" {
		r.Handle("/*", spaHandler(s.cfg.StaticDir))
	}

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, httpx.ErrNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, &httpx.APIError{Code: "method_not_allowed", Message: "method not allowed"})
	})

	return r
}

func (s *Server) requestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			slog.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"dur_ms", time.Since(start).Milliseconds(),
				"req_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
