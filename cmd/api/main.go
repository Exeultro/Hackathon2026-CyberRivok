package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"typing/internal/auth"
	"typing/internal/config"
	"typing/internal/gamification"
	"typing/internal/httpserver"
	"typing/internal/leaderboard"
	"typing/internal/maps"
	"typing/internal/progress"
	"typing/internal/sessions"
	"typing/internal/stats"
	"typing/internal/store/postgres"
	"typing/internal/texts"
	"typing/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, cfg.DB)
	if err != nil {
		slog.Error("connect db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := runMigrations(ctx, pool); err != nil {
		slog.Error("migrate", "err", err)
		os.Exit(1)
	}

	if err := gamification.Seed(ctx, pool); err != nil {
		slog.Error("seed gamification", "err", err)
		os.Exit(1)
	}
	if err := texts.Seed(ctx, pool); err != nil {
		slog.Error("seed texts", "err", err)
		os.Exit(1)
	}

	authRepo := auth.NewRepo(pool)
	authSvc := auth.NewService(authRepo, cfg)
	authH := auth.NewHandler(authSvc)

	userRepo := user.NewRepo(pool)
	userSvc := user.NewService(userRepo)
	userH := user.NewHandler(userSvc)

	progressRepo := progress.NewRepo(pool)
	progressSvc := progress.NewService(progressRepo)
	progressH := progress.NewHandler(progressSvc)

	mapsRepo := maps.NewRepo(pool)
	mapsSvc := maps.NewService(mapsRepo)
	mapsH := maps.NewHandler(mapsSvc, progressSvc)

	sessionsRepo := sessions.NewRepo(pool)
	sessionsSvc := sessions.NewService(sessionsRepo, progressSvc, mapsSvc)
	sessionsH := sessions.NewHandler(sessionsSvc)

	statsRepo := stats.NewRepo(pool)
	statsH := stats.NewHandler(statsRepo)

	lbRepo := leaderboard.NewRepo(pool)
	lbH := leaderboard.NewHandler(lbRepo)

	textsRepo := texts.NewRepo(pool)
	textsSvc := texts.NewService(textsRepo)
	textsH := texts.NewHandler(textsSvc)

	srv := httpserver.New(cfg, authH, userH, sessionsH, progressH, statsH, lbH, textsH, mapsH)
	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      srv.Router(),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		slog.Info("http server starting", "addr", cfg.HTTP.Addr, "env", cfg.Env)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "err", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	goose.SetBaseFS(nil) // миграции с диска
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	return goose.UpContext(ctx, db, "internal/store/migrations")
}
