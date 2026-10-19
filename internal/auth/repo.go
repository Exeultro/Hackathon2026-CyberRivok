package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists      = errors.New("auth: user already exists")
	ErrUserNotFound    = errors.New("auth: user not found")
	ErrRefreshNotFound = errors.New("auth: refresh token not found")
	ErrEmailTaken      = errors.New("auth: email taken")
	ErrNicknameTaken   = errors.New("auth: nickname taken")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Nickname     string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash []byte
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) CreateUser(ctx context.Context, email, nickname, passwordHash string) (*User, error) {
	const q = `
		INSERT INTO users (email, nickname, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, nickname, password_hash, created_at`

	var u User
	err := r.pool.QueryRow(ctx, q, email, nickname, passwordHash).
		Scan(&u.ID, &u.Email, &u.Nickname, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, ErrEmailTaken
			case "users_nickname_key":
				return nil, ErrNicknameTaken
			default:
				return nil, ErrUserExists
			}
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const q = `SELECT id, email, nickname, password_hash, created_at FROM users WHERE email = $1`
	var u User
	err := r.pool.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.Nickname, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *Repo) GetUserByID(ctx context.Context, id string) (*User, error) {
	const q = `SELECT id, email, nickname, password_hash, created_at FROM users WHERE id = $1`
	var u User
	err := r.pool.QueryRow(ctx, q, id).
		Scan(&u.ID, &u.Email, &u.Nickname, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (r *Repo) InsertRefreshToken(ctx context.Context, userID string, hash []byte, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, q, userID, hash, expiresAt)
	return err
}

func (r *Repo) GetActiveRefresh(ctx context.Context, hash []byte) (*RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()`
	var rt RefreshToken
	err := r.pool.QueryRow(ctx, q, hash).
		Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRefreshNotFound
	}
	return &rt, err
}

func (r *Repo) RevokeRefresh(ctx context.Context, hash []byte) error {
	const q = `UPDATE refresh_tokens SET revoked_at = now() WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, q, hash)
	return err
}
