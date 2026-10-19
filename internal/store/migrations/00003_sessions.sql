-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode                text NOT NULL CHECK (mode IN ('words','sentences','code','symbols')),
    lang                text NOT NULL CHECK (lang IN ('ru','en')),
    code_language       text NULL,
    difficulty          int  NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    started_at          timestamptz NOT NULL,
    duration_ms         int  NOT NULL CHECK (duration_ms > 0),
    chars_typed         int  NOT NULL CHECK (chars_typed >= 0),
    chars_correct       int  NOT NULL CHECK (chars_correct >= 0),
    errors              int  NOT NULL CHECK (errors >= 0),
    max_combo           int  NOT NULL DEFAULT 0,
    longest_clean_streak int NOT NULL DEFAULT 0,
    client_score        int  NULL,
    wpm                 numeric(8,2) NOT NULL,
    cpm                 numeric(8,2) NOT NULL,
    accuracy            numeric(5,2) NOT NULL,
    score               int NOT NULL,
    xp_gained           int NOT NULL,
    created_at          timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_sessions_user_started ON sessions(user_id, started_at DESC);
CREATE INDEX idx_sessions_started ON sessions(started_at DESC);
CREATE INDEX idx_sessions_leaderboard ON sessions(started_at DESC, mode, lang, wpm DESC);

CREATE TABLE session_key_stats (
    id             bigserial PRIMARY KEY,
    session_id     uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id        uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind           text NOT NULL CHECK (kind IN ('key','bigram')),
    key            text NOT NULL,
    lang           text NOT NULL CHECK (lang IN ('ru','en')),
    hits           int  NOT NULL CHECK (hits >= 0),
    errors         int  NOT NULL CHECK (errors >= 0),
    avg_latency_ms numeric(10,2) NULL,
    created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_sks_user_kind_lang ON session_key_stats(user_id, kind, lang);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session_key_stats;
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd