-- +goose Up
-- +goose StatementBegin
CREATE TYPE server_lifecycle_state AS ENUM (
    'active',
    'pending_deletion'
);

ALTER TABLE servers
    ADD COLUMN lifecycle_state server_lifecycle_state
    NOT NULL DEFAULT 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE servers
    DROP COLUMN lifecycle_state;

DROP TYPE IF EXISTS server_lifecycle_state;
-- +goose StatementEnd