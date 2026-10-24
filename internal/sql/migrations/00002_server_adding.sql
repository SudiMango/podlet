-- +goose Up
-- +goose StatementBegin
CREATE TABLE pending_servers (
    id VARCHAR(64) NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    registration_token_hash VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE servers (
    id VARCHAR(64) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT name_user_id_unique
        UNIQUE (name, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pending_servers;
DROP TABLE IF EXISTS servers;
-- +goose StatementEnd
