package gamification

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for lvl := 1; lvl <= 100; lvl++ {
		if _, err := tx.Exec(ctx,
			`INSERT INTO levels(level, xp_required) VALUES ($1,$2)
			 ON CONFLICT (level) DO UPDATE SET xp_required = EXCLUDED.xp_required`,
			lvl, LevelXPThreshold(lvl),
		); err != nil {
			return err
		}
	}

	for _, u := range Unlocks() {
		if _, err := tx.Exec(ctx,
			`INSERT INTO unlocks(level, type, value) VALUES ($1,$2,$3)
			 ON CONFLICT DO NOTHING`,
			u.Level, u.Type, u.Value,
		); err != nil {
			return err
		}
	}

	for _, a := range Achievements() {
		if _, err := tx.Exec(ctx,
			`INSERT INTO achievements(code,title,description,icon,target,metric)
			 VALUES ($1,$2,$3,$4,$5,$6)
			 ON CONFLICT (code) DO UPDATE
			   SET title=EXCLUDED.title, description=EXCLUDED.description,
			       icon=EXCLUDED.icon, target=EXCLUDED.target, metric=EXCLUDED.metric`,
			a.Code, a.Title, a.Description, a.Icon, a.Target, a.Metric,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
