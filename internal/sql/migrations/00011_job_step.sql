-- +goose Up
-- +goose StatementBegin
ALTER TABLE server_agent_jobs ADD COLUMN step TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE server_agent_jobs DROP COLUMN IF EXISTS step;
-- +goose StatementEnd
