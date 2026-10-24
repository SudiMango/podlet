-- +goose Up
-- +goose StatementBegin
ALTER TABLE servers
    ADD COLUMN last_seen_at TIMESTAMPTZ DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE servers
    DROP COLUMN last_seen_at;
-- +goose StatementEnd
