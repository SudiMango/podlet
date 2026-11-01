-- +goose Up
-- +goose StatementBegin
ALTER TABLE server_agent_jobs
    ALTER COLUMN job_type TYPE TEXT
    USING job_type::TEXT;

DROP TYPE server_job_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE server_job_type AS ENUM (
    'clone_repository'
);

ALTER TABLE server_agent_jobs
    ALTER COLUMN job_type TYPE server_job_type
    USING job_type::server_job_type;
-- +goose StatementEnd
