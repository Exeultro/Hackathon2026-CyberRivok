package sessions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("sessions: not found")

type Session struct {
	ID                 string
	UserID             string
	Mode               string
	Lang               string
	CodeLanguage       *string
	Difficulty         int
	StartedAt          time.Time
	DurationMS         int
	CharsTyped         int
	CharsCorrect       int
	Errors             int
	MaxCombo           int
	LongestCleanStreak int
	ClientScore        *int
	WPM                float64
	CPM                float64
	Accuracy           float64
	Score              int
	XPGained           int
	CreatedAt          time.Time
}

type KeyStatRow struct {
	Kind         string
	Key          string
	Lang         string
	Hits         int
	Errors       int
	AvgLatencyMS *float64
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type CreateInput struct {
	Session  Session
	KeyStats []KeyStatRow
}

func (r *Repo) Create(ctx context.Context, in CreateInput) (*Session, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	s := in.Session
	const q = `
		INSERT INTO sessions (
			id, user_id, mode, lang, code_language, difficulty,
			started_at, duration_ms, chars_typed, chars_correct, errors,
			max_combo, longest_clean_streak, client_score,
			wpm, cpm, accuracy, score, xp_gained
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13,
			$14, $15, $16, $17, $18
		) RETURNING id, created_at`

	err = tx.QueryRow(ctx, q,
		s.UserID, s.Mode, s.Lang, s.CodeLanguage, s.Difficulty,
		s.StartedAt, s.DurationMS, s.CharsTyped, s.CharsCorrect, s.Errors,
		s.MaxCombo, s.LongestCleanStreak, s.ClientScore,
		s.WPM, s.CPM, s.Accuracy, s.Score, s.XPGained,
	).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}

	const kq = `
		INSERT INTO session_key_stats (session_id, user_id, kind, key, lang, hits, errors, avg_latency_ms)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	for _, k := range in.KeyStats {
		if _, err := tx.Exec(ctx, kq, s.ID, in.Session.UserID, k.Kind, k.Key, k.Lang, k.Hits, k.Errors, k.AvgLatencyMS); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	in.Session.ID = s.ID
	in.Session.CreatedAt = s.CreatedAt
	return &in.Session, nil
}

func (r *Repo) List(ctx context.Context, userID string, limit, offset int) ([]Session, int, error) {
	const cq = `SELECT count(*) FROM sessions WHERE user_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, cq, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `
		SELECT id, user_id, mode, lang, code_language, difficulty,
		       started_at, duration_ms, chars_typed, chars_correct, errors,
		       max_combo, longest_clean_streak, client_score,
		       wpm, cpm, accuracy, score, xp_gained, created_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.Mode, &s.Lang, &s.CodeLanguage, &s.Difficulty,
			&s.StartedAt, &s.DurationMS, &s.CharsTyped, &s.CharsCorrect, &s.Errors,
			&s.MaxCombo, &s.LongestCleanStreak, &s.ClientScore,
			&s.WPM, &s.CPM, &s.Accuracy, &s.Score, &s.XPGained, &s.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

func (r *Repo) BestWPMBefore(ctx context.Context, userID string, before time.Time) (float64, error) {
	const q = `SELECT COALESCE(MAX(wpm), 0) FROM sessions WHERE user_id = $1 AND started_at < $2`
	var v float64
	err := r.pool.QueryRow(ctx, q, userID, before).Scan(&v)
	return v, err
}

func (r *Repo) MaxCleanStreak(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COALESCE(MAX(longest_clean_streak), 0) FROM sessions WHERE user_id = $1`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) TotalChars(ctx context.Context, userID string) (int64, error) {
	const q = `SELECT COALESCE(SUM(chars_correct), 0) FROM sessions WHERE user_id = $1`
	var v int64
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) SessionsCount(ctx context.Context, userID string) (int, error) {
	const q = `SELECT count(*) FROM sessions WHERE user_id = $1`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) BestCPM(ctx context.Context, userID string) (float64, error) {
	const q = `SELECT COALESCE(MAX(cpm), 0) FROM sessions WHERE user_id = $1`
	var v float64
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) MaxCombo(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COALESCE(MAX(max_combo), 0) FROM sessions WHERE user_id = $1`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) CodeSessionsCount(ctx context.Context, userID string) (int, error) {
	const q = `SELECT count(*) FROM sessions WHERE user_id = $1 AND mode = 'code'`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) DistinctLangs(ctx context.Context, userID string) (int, error) {
	const q = `SELECT count(DISTINCT lang) FROM sessions WHERE user_id = $1`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) PerfectSessionsCount(ctx context.Context, userID string) (int, error) {
	const q = `SELECT count(*) FROM sessions WHERE user_id = $1 AND accuracy = 100 AND chars_typed >= 20`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}

func (r *Repo) StreakDays(ctx context.Context, userID string) (int, error) {
	const q = `
		SELECT DISTINCT (started_at AT TIME ZONE 'UTC')::date AS d
		FROM sessions WHERE user_id = $1 ORDER BY d DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var days []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return 0, err
		}
		days = append(days, d)
	}
	if len(days) == 0 {
		return 0, nil
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	var expected time.Time
	switch {
	case days[0].Equal(today):
		expected = today
	case days[0].Equal(yesterday):
		expected = yesterday
	default:
		return 0, nil
	}
	streak := 0
	for _, d := range days {
		if d.Equal(expected) {
			streak++
			expected = expected.AddDate(0, 0, -1)
		} else if d.Before(expected) {
			break
		}
	}
	return streak, nil
}

func (r *Repo) LangsWithCurrent(ctx context.Context, userID, currentLang string) (int, error) {
	const q = `
		SELECT count(DISTINCT lang) FROM (
			SELECT lang FROM sessions WHERE user_id = $1
			UNION
			SELECT $2::text
		) t`
	var v int
	err := r.pool.QueryRow(ctx, q, userID, currentLang).Scan(&v)
	return v, err
}
