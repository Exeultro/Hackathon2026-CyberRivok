-- +goose Up
-- +goose StatementBegin
CREATE TABLE texts (
    id            text PRIMARY KEY,
    mode          text NOT NULL CHECK (mode IN ('sentences','code')),
    lang          text NOT NULL CHECK (lang IN ('ru','en')),
    difficulty    int  NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
    code_language text NULL CHECK (code_language IN ('python','javascript','html','css')),
    text          text NOT NULL
);
CREATE INDEX idx_texts_lookup ON texts(mode, lang, difficulty, code_language);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS texts;
-- +goose StatementEnd