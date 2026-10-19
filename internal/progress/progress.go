package progress

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"typing/internal/gamification"
)

var ErrNotFound = errors.New("progress: not found")

type Progress struct {
	Level          int      `json:"level"`
	XPTotal        int64    `json:"xp_total"`
	XPInLevel      int64    `json:"xp_in_level"`
	XPToNext       int64    `json:"xp_to_next"`
	UnlockedModes  []string `json:"unlocked_modes"`
	MaxDifficulty  int      `json:"max_difficulty"`
	UnlockedThemes []string `json:"unlocked_themes,omitempty"`
}

type Achievement struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Target      float64 `json:"target"`
	Progress    float64 `json:"progress"`
	UnlockedAt  *string `json:"unlocked_at,omitempty"`
}

type Unlock struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ApplyResult struct {
	LevelBefore     int
	Progress        Progress
	NewAchievements []Achievement
	Unlocked        []Unlock
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// EnsureRow — создаёт user_progress, если нет.
func (r *Repo) EnsureRow(ctx context.Context, userID string) error {
	const q = `INSERT INTO user_progress(user_id) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}

func (r *Repo) GetXP(ctx context.Context, userID string) (int64, int, error) {
	const q = `SELECT xp_total, level FROM user_progress WHERE user_id = $1`
	var xp int64
	var lvl int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&xp, &lvl)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 1, nil
	}
	return xp, lvl, err
}

func (r *Repo) SetXPAndLevel(ctx context.Context, userID string, xp int64, level int) error {
	const q = `UPDATE user_progress SET xp_total=$2, level=$3, updated_at=now() WHERE user_id=$1`
	_, err := r.pool.Exec(ctx, q, userID, xp, level)
	return err
}

func (r *Repo) UnlocksUpTo(ctx context.Context, level int) ([]Unlock, error) {
	const q = `SELECT type, value FROM unlocks WHERE level <= $1 ORDER BY level`
	rows, err := r.pool.Query(ctx, q, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Unlock
	for rows.Next() {
		var u Unlock
		if err := rows.Scan(&u.Type, &u.Value); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *Repo) UnlocksInRange(ctx context.Context, fromLevel, toLevel int) ([]Unlock, error) {
	if toLevel <= fromLevel {
		return nil, nil
	}
	const q = `SELECT type, value FROM unlocks WHERE level > $1 AND level <= $2 ORDER BY level`
	rows, err := r.pool.Query(ctx, q, fromLevel, toLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Unlock
	for rows.Next() {
		var u Unlock
		if err := rows.Scan(&u.Type, &u.Value); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// BuildProgress собирает Progress из xp_total.
func (r *Repo) BuildProgress(ctx context.Context, userID string, xp int64) (Progress, error) {
	lvl := gamification.LevelFromXP(xp)
	unlocks, err := r.UnlocksUpTo(ctx, lvl)
	if err != nil {
		return Progress{}, err
	}
	p := Progress{
		Level:     lvl,
		XPTotal:   xp,
		XPInLevel: gamification.XPInLevel(xp),
		XPToNext:  gamification.XPToNext(xp),
	}
	for _, u := range unlocks {
		switch u.Type {
		case "mode":
			p.UnlockedModes = append(p.UnlockedModes, u.Value)
		case "difficulty":
			if d := atoi(u.Value); d > p.MaxDifficulty {
				p.MaxDifficulty = d
			}
		case "theme":
			p.UnlockedThemes = append(p.UnlockedThemes, u.Value)
		}
	}
	return p, nil
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func (r *Repo) CatalogAll(ctx context.Context) ([]Achievement, error) {
	const q = `SELECT code, title, description, icon, target FROM achievements ORDER BY code`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Achievement, 0)
	for rows.Next() {
		var a Achievement
		if err := rows.Scan(&a.Code, &a.Title, &a.Description, &a.Icon, &a.Target); err != nil {
			return nil, err
		}
		a.Progress = 0
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) AchievementProgress(ctx context.Context, userID string) (map[string]Achievement, error) {
	const q = `
		SELECT a.code, a.title, a.description, a.icon, a.target,
		       COALESCE(ua.progress, 0), ua.unlocked_at
		FROM achievements a
		LEFT JOIN user_achievements ua ON ua.code = a.code AND ua.user_id = $1`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]Achievement{}
	for rows.Next() {
		var a Achievement
		var unlockedAt *time.Time
		if err := rows.Scan(&a.Code, &a.Title, &a.Description, &a.Icon, &a.Target, &a.Progress, &unlockedAt); err != nil {
			return nil, err
		}
		if unlockedAt != nil {
			s := unlockedAt.UTC().Format(time.RFC3339)
			a.UnlockedAt = &s
		}
		out[a.Code] = a
	}
	return out, rows.Err()
}

func (r *Repo) UpsertAchievement(ctx context.Context, userID, code string, progress float64, unlocked bool) error {
	const q = `
		INSERT INTO user_achievements(user_id, code, progress, unlocked_at)
		VALUES ($1,$2,$3, CASE WHEN $4 THEN now() ELSE NULL END)
		ON CONFLICT (user_id, code) DO UPDATE
		SET progress = EXCLUDED.progress,
		    unlocked_at = COALESCE(user_achievements.unlocked_at, EXCLUDED.unlocked_at)`
	_, err := r.pool.Exec(ctx, q, userID, code, progress, unlocked)
	return err
}

func (r *Repo) Reset(ctx context.Context, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM session_key_stats WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM user_achievements WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM user_progress WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
