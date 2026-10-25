-- +goose Up
-- +goose StatementBegin
CREATE TABLE github_installation_states (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL 
);

CREATE TABLE github_installations (
    id BIGINT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    account_id BIGINT NOT NULL,
    account_login TEXT NOT NULL,
    account_type TEXT NOT NULL,
    repository_selection TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS github_installation_states;
DROP TABLE IF EXISTS github_installations;
-- +goose StatementEnd