package leaderboard

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Entry struct {
	Rank     int     `json:"rank"`
	UserID   string  `json:"-"`
	Nickname string  `json:"-"`
	Level    int     `json:"-"`
	Value    float64 `json:"value"`
}

type EntryUser struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Level    int    `json:"level"`
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type Filter struct {
	Metric string // speed | accuracy | xp
	Period string // today | week | all
	Mode   string
	Lang   string
	Limit  int
	Offset int
}

func (r *Repo) List(ctx context.Context, f Filter) ([]Entry, int, error) {
	from := periodFrom(f.Period)

	var (
		q    string
		args []any
	)
	switch f.Metric {
	case "speed":
		q = `
			SELECT s.user_id, u.nickname, COALESCE(up.level,1) AS level, max(s.wpm)::float8 AS value
			FROM sessions s
			JOIN users u ON u.id = s.user_id
			LEFT JOIN user_progress up ON up.user_id = s.user_id
			WHERE ($1::timestamptz IS NULL OR s.started_at >= $1)
			  AND ($2 = '' OR s.mode = $2)
			  AND ($3 = '' OR s.lang = $3)
			GROUP BY s.user_id, u.nickname, up.level
			ORDER BY value DESC, s.user_id
			LIMIT $4 OFFSET $5`
		args = []any{from, f.Mode, f.Lang, f.Limit, f.Offset}

	case "accuracy":
		q = `
			SELECT s.user_id, u.nickname, COALESCE(up.level,1) AS level, avg(s.accuracy)::float8 AS value
			FROM sessions s
			JOIN users u ON u.id = s.user_id
			LEFT JOIN user_progress up ON up.user_id = s.user_id
			WHERE ($1::timestamptz IS NULL OR s.started_at >= $1)
			  AND ($2 = '' OR s.mode = $2)
			  AND ($3 = '' OR s.lang = $3)
			GROUP BY s.user_id, u.nickname, up.level
			HAVING count(*) >= 3
			ORDER BY value DESC, s.user_id
			LIMIT $4 OFFSET $5`
		args = []any{from, f.Mode, f.Lang, f.Limit, f.Offset}

	case "xp":
		q = `
			SELECT s.user_id, u.nickname, COALESCE(up.level,1) AS level, sum(s.xp_gained)::float8 AS value
			FROM sessions s
			JOIN users u ON u.id = s.user_id
			LEFT JOIN user_progress up ON up.user_id = s.user_id
			WHERE ($1::timestamptz IS NULL OR s.started_at >= $1)
			  AND ($2 = '' OR s.mode = $2)
			  AND ($3 = '' OR s.lang = $3)
			GROUP BY s.user_id, u.nickname, up.level
			ORDER BY value DESC, s.user_id
			LIMIT $4 OFFSET $5`
		args = []any{from, f.Mode, f.Lang, f.Limit, f.Offset}

	default:
		return nil, 0, nil
	}

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]Entry, 0)
	rank := f.Offset
	for rows.Next() {
		rank++
		var e Entry
		if err := rows.Scan(&e.UserID, &e.Nickname, &e.Level, &e.Value); err != nil {
			return nil, 0, err
		}
		e.Rank = rank
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.total(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *Repo) total(ctx context.Context, f Filter) (int, error) {
	from := periodFrom(f.Period)
	var q string
	switch f.Metric {
	case "accuracy":
		q = `
			SELECT count(*) FROM (
				SELECT s.user_id FROM sessions s
				WHERE ($1::timestamptz IS NULL OR s.started_at >= $1)
				  AND ($2 = '' OR s.mode = $2)
				  AND ($3 = '' OR s.lang = $3)
				GROUP BY s.user_id HAVING count(*) >= 3
			) t`
	default:
		q = `
			SELECT count(DISTINCT s.user_id) FROM sessions s
			WHERE ($1::timestamptz IS NULL OR s.started_at >= $1)
			  AND ($2 = '' OR s.mode = $2)
			  AND ($3 = '' OR s.lang = $3)`
	}
	var n int
	err := r.pool.QueryRow(ctx, q, from, f.Mode, f.Lang).Scan(&n)
	return n, err
}

func (r *Repo) Me(ctx context.Context, f Filter, userID string) (*Entry, error) {
	f2 := f
	f2.Limit = 1000
	f2.Offset = 0
	items, _, err := r.List(ctx, f2)
	if err != nil {
		return nil, err
	}
	for _, e := range items {
		if e.UserID == userID {
			return &e, nil
		}
	}
	return nil, nil
}

func periodFrom(period string) *time.Time {
	now := time.Now().UTC()
	switch period {
	case "today":
		t := now.Truncate(24 * time.Hour)
		return &t
	case "week":
		t := now.AddDate(0, 0, -7)
		return &t
	case "all", "":
		return nil
	}
	return nil
}
