-- +goose Up
-- +goose StatementBegin
ALTER TABLE texts ADD COLUMN source text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE texts DROP COLUMN IF EXISTS source;
-- +goose StatementEnd