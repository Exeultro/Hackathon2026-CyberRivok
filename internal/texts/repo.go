package texts

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID           string  `json:"id"`
	Mode         string  `json:"mode"`
	Lang         string  `json:"lang"`
	Difficulty   int     `json:"difficulty"`
	CodeLanguage *string `json:"code_language,omitempty"`
	Text         string  `json:"text"`
	Source       *string `json:"source,omitempty"`
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Pick(ctx context.Context, mode, lang string, difficulty, count int, codeLanguage string) ([]Item, error) {
	const q = `
		SELECT id, mode, lang, difficulty, code_language, text, source
		FROM texts
		WHERE mode = $1 AND lang = $2 AND difficulty = $3
		  AND ($4 = '' OR code_language = $4)
		ORDER BY random()
		LIMIT $5`

	rows, err := r.pool.Query(ctx, q, mode, lang, difficulty, codeLanguage, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Item, 0, count)
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Mode, &it.Lang, &it.Difficulty, &it.CodeLanguage, &it.Text, &it.Source); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
