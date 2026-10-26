-- +goose Up
-- +goose StatementBegin
CREATE TYPE server_job_status AS ENUM (
    'pending',
    'ongoing',
    'success',
    'failed'
);

CREATE TYPE server_job_type AS ENUM (
    'clone_repository'
);

CREATE TABLE server_agent_jobs (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    server_id VARCHAR(64) NOT NULL REFERENCES servers (id),
    job_status server_job_status NOT NULL DEFAULT 'pending',
    job_type server_job_type NOT NULL,
    repo_name TEXT NOT NULL,
    repository_id BIGINT NOT NULL,
    installation_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS server_agent_jobs;

DROP TYPE IF EXISTS server_job_type;
DROP TYPE IF EXISTS server_job_status;
-- +goose StatementEnd