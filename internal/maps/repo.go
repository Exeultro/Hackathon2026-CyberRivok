package maps

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeResult struct {
	NodeID       string
	Stars        int
	BestWPM      float64
	BestAccuracy float64
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) ResultsForUser(ctx context.Context, userID string) (map[string]NodeResult, error) {
	const q = `SELECT node_id, stars, best_wpm, best_accuracy FROM map_node_results WHERE user_id = $1`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]NodeResult{}
	for rows.Next() {
		var n NodeResult
		if err := rows.Scan(&n.NodeID, &n.Stars, &n.BestWPM, &n.BestAccuracy); err != nil {
			return nil, err
		}
		out[n.NodeID] = n
	}
	return out, rows.Err()
}

func (r *Repo) UpsertBest(ctx context.Context, userID, nodeID string, stars int, wpm, accuracy float64) (NodeResult, error) {
	const q = `
		INSERT INTO map_node_results(user_id, node_id, stars, best_wpm, best_accuracy)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (user_id, node_id) DO UPDATE
		SET stars = GREATEST(map_node_results.stars, EXCLUDED.stars),
		    best_wpm = CASE
		        WHEN EXCLUDED.stars > map_node_results.stars THEN EXCLUDED.best_wpm
		        WHEN EXCLUDED.stars = map_node_results.stars AND EXCLUDED.best_wpm > map_node_results.best_wpm THEN EXCLUDED.best_wpm
		        ELSE map_node_results.best_wpm
		    END,
		    best_accuracy = CASE
		        WHEN EXCLUDED.stars > map_node_results.stars THEN EXCLUDED.best_accuracy
		        WHEN EXCLUDED.stars = map_node_results.stars AND EXCLUDED.best_wpm = map_node_results.best_wpm AND EXCLUDED.best_accuracy > map_node_results.best_accuracy THEN EXCLUDED.best_accuracy
		        ELSE map_node_results.best_accuracy
		    END,
		    updated_at = now()
		RETURNING node_id, stars, best_wpm, best_accuracy`
	var n NodeResult
	err := r.pool.QueryRow(ctx, q, userID, nodeID, stars, wpm, accuracy).
		Scan(&n.NodeID, &n.Stars, &n.BestWPM, &n.BestAccuracy)
	return n, err
}

func (r *Repo) StarsTotal(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COALESCE(SUM(stars), 0) FROM map_node_results WHERE user_id = $1`
	var v int
	err := r.pool.QueryRow(ctx, q, userID).Scan(&v)
	return v, err
}
