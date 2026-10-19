package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"typing/internal/settings"
)

var (
	ErrNotFound      = errors.New("user: not found")
	ErrNicknameTaken = errors.New("user: nickname taken")
)

type User struct {
	ID        string
	Email     string
	Nickname  string
	CreatedAt string
	Settings  settings.Settings
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) GetByID(ctx context.Context, id string) (*User, error) {
	const q = `
		SELECT id, email, nickname,
		       to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at,
		       settings
		FROM users WHERE id = $1`

	var (
		u      User
		rawSet []byte
	)
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Email, &u.Nickname, &u.CreatedAt, &rawSet)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	u.Settings = settings.Default()
	if len(rawSet) > 0 {
		var patch settings.Patch
		if err := json.Unmarshal(rawSet, &patch); err != nil {
			return nil, err
		}
		u.Settings = settings.Apply(u.Settings, patch)
	}
	return &u, nil
}

func (r *Repo) Update(ctx context.Context, id string, nickname *string, patch *settings.Patch) (*User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if nickname != nil {
		const q = `UPDATE users SET nickname = $1, updated_at = now() WHERE id = $2`
		if _, err := tx.Exec(ctx, q, *nickname, id); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_nickname_key" {
				return nil, ErrNicknameTaken
			}
			return nil, err
		}
	}

	if patch != nil {
		const q = `
			UPDATE users
			SET settings = COALESCE(settings, '{}'::jsonb) || $1::jsonb,
			    updated_at = now()
			WHERE id = $2`
		patchBytes, err := json.Marshal(patch)
		if err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, q, patchBytes, id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}
