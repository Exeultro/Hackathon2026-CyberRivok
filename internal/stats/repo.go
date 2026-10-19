package stats

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Summary struct {
	SessionsCount int     `json:"sessions_count"`
	TotalTimeMS   int64   `json:"total_time_ms"`
	TotalChars    int64   `json:"total_chars"`
	BestWPM       float64 `json:"best_wpm"`
	BestCPM       float64 `json:"best_cpm"`
	AvgWPM        float64 `json:"avg_wpm"`
	AvgAccuracy   float64 `json:"avg_accuracy"`
	StreakDays    int     `json:"streak_days"`
}

type TimelinePoint struct {
	Date        string  `json:"date"`
	Sessions    int     `json:"sessions"`
	AvgWPM      float64 `json:"avg_wpm"`
	AvgCPM      float64 `json:"avg_cpm"`
	AvgAccuracy float64 `json:"avg_accuracy"`
	Errors      int     `json:"errors,omitempty"`
}

type KeyStat struct {
	Key          string   `json:"key"`
	Hits         int      `json:"hits"`
	Errors       int      `json:"errors"`
	ErrorRate    float64  `json:"error_rate"`
	AvgLatencyMS *float64 `json:"avg_latency_ms,omitempty"`
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Summary(ctx context.Context, userID string) (Summary, error) {
	const q = `
		SELECT
			COALESCE(count(*), 0)                              AS sessions_count,
			COALESCE(sum(duration_ms), 0)                      AS total_time_ms,
			COALESCE(sum(chars_correct), 0)                    AS total_chars,
			COALESCE(max(wpm), 0)                              AS best_wpm,
			COALESCE(max(cpm), 0)                              AS best_cpm,
			COALESCE(avg(wpm), 0)                              AS avg_wpm,
			COALESCE(avg(accuracy), 0)                         AS avg_accuracy
		FROM sessions WHERE user_id = $1`

	var s Summary
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&s.SessionsCount, &s.TotalTimeMS, &s.TotalChars,
		&s.BestWPM, &s.BestCPM, &s.AvgWPM, &s.AvgAccuracy,
	)
	if err != nil {
		return Summary{}, err
	}

	streak, err := r.streakDays(ctx, userID)
	if err != nil {
		return Summary{}, err
	}
	s.StreakDays = streak
	return s, nil
}

func (r *Repo) streakDays(ctx context.Context, userID string) (int, error) {
	const q = `
		SELECT DISTINCT (started_at AT TIME ZONE 'UTC')::date AS d
		FROM sessions
		WHERE user_id = $1
		ORDER BY d DESC`
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

func (r *Repo) Timeline(ctx context.Context, userID string, from time.Time) ([]TimelinePoint, error) {
	const q = `
		SELECT
			to_char((started_at AT TIME ZONE 'UTC')::date, 'YYYY-MM-DD') AS d,
			count(*)                        AS sessions,
			COALESCE(avg(wpm), 0)           AS avg_wpm,
			COALESCE(avg(cpm), 0)           AS avg_cpm,
			COALESCE(avg(accuracy), 0)      AS avg_accuracy,
			COALESCE(sum(errors), 0)        AS errors
		FROM sessions
		WHERE user_id = $1 AND started_at >= $2
		GROUP BY d
		ORDER BY d`

	rows, err := r.pool.Query(ctx, q, userID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]TimelinePoint, 0)
	for rows.Next() {
		var p TimelinePoint
		if err := rows.Scan(&p.Date, &p.Sessions, &p.AvgWPM, &p.AvgCPM, &p.AvgAccuracy, &p.Errors); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repo) KeyStats(ctx context.Context, userID, kind, lang string, limit int) ([]KeyStat, error) {
	const q = `
		SELECT key,
		       sum(hits)::int   AS hits,
		       sum(errors)::int AS errors,
		       avg(avg_latency_ms)::float8 AS avg_latency_ms
		FROM session_key_stats
		WHERE user_id = $1 AND kind = $2 AND ($3 = '' OR lang = $3)
		GROUP BY key
		HAVING sum(hits) >= 3 AND sum(errors) > 0
		ORDER BY (sum(errors)::float8 / NULLIF(sum(hits),0)) DESC, hits DESC
		LIMIT $4`

	rows, err := r.pool.Query(ctx, q, userID, kind, lang, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]KeyStat, 0)
	for rows.Next() {
		var k KeyStat
		if err := rows.Scan(&k.Key, &k.Hits, &k.Errors, &k.AvgLatencyMS); err != nil {
			return nil, err
		}
		if k.Hits > 0 {
			k.ErrorRate = float64(k.Errors) / float64(k.Hits)
		}
		out = append(out, k)
	}
	return out, rows.Err()
}
