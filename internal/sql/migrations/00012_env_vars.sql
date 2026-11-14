-- +goose Up
-- +goose StatementBegin
CREATE TABLE env_vars (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    key VARCHAR(64) NOT NULL,
    value_encrypted TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    repository_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_id, repository_id, key)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS env_vars;
-- +goose StatementEnd
