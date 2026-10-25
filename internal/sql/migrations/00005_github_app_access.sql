-- +goose Up
-- +goose StatementBegin
CREATE TABLE github_installation_states (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL 
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS github_installation_states
-- +goose StatementEnd