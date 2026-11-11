-- +goose Up
-- +goose StatementBegin
ALTER TABLE dns_records
    ADD COLUMN repository_id BIGINT NOT NULL,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE UNIQUE INDEX dns_records_server_repository_idx
    ON dns_records (server_id, repository_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS dns_records_server_repository_idx;

ALTER TABLE dns_records
    DROP COLUMN IF EXISTS repository_id,
    DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
