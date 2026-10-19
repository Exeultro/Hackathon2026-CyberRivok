-- +goose Up
-- +goose StatementBegin
CREATE TABLE map_node_results (
    user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    node_id        text NOT NULL,
    stars          int  NOT NULL DEFAULT 0 CHECK (stars BETWEEN 0 AND 3),
    best_wpm       numeric(8,2) NOT NULL DEFAULT 0,
    best_accuracy  numeric(5,2) NOT NULL DEFAULT 0,
    updated_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, node_id)
);
CREATE INDEX idx_mnr_user ON map_node_results(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS map_node_results;
-- +goose StatementEnd