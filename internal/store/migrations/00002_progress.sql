-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_progress (
    user_id     uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    xp_total    bigint NOT NULL DEFAULT 0,
    level       int    NOT NULL DEFAULT 1,
    updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_user_progress_level ON user_progress(level);
CREATE INDEX idx_user_progress_xp ON user_progress(xp_total DESC);

CREATE TABLE levels (
    level         int PRIMARY KEY,
    xp_required   bigint NOT NULL  -- сколько XP суммарно нужно, чтобы ДОСТИЧЬ этот уровень
);

CREATE TABLE unlocks (
    level int  NOT NULL REFERENCES levels(level) ON DELETE CASCADE,
    type  text NOT NULL CHECK (type IN ('mode','difficulty','theme')),
    value text NOT NULL,
    PRIMARY KEY (level, type, value)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS unlocks;
DROP TABLE IF EXISTS levels;
DROP TABLE IF EXISTS user_progress;
-- +goose StatementEnd