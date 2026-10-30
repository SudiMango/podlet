-- +goose Up
-- +goose StatementBegin
CREATE TABLE tunnels (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    name VARCHAR(64) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    server_id VARCHAR(64) UNIQUE NOT NULL REFERENCES servers (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dns_records (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    name VARCHAR(64) UNIQUE NOT NULL,
    service VARCHAR(64) NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id),
    server_id VARCHAR(64) NOT NULL REFERENCES servers (id),
    tunnel_id VARCHAR(64) NOT NULL REFERENCES tunnels (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dns_records;
DROP TABLE IF EXISTS tunnels;
-- +goose StatementEnd
