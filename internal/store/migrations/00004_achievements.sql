-- +goose Up
-- +goose StatementBegin
CREATE TABLE achievements (
    code        text PRIMARY KEY,
    title       text NOT NULL,
    description text NOT NULL,
    icon        text NOT NULL,
    target      numeric(12,2) NOT NULL,
    metric      text NOT NULL  -- напр. best_wpm, clean_streak, total_chars, sessions_count
);

CREATE TABLE user_achievements (
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code        text NOT NULL REFERENCES achievements(code) ON DELETE CASCADE,
    progress    numeric(12,2) NOT NULL DEFAULT 0,
    unlocked_at timestamptz NULL,
    PRIMARY KEY (user_id, code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievements;
-- +goose StatementEnd