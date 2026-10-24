-- name: CreatePendingServerByUserId :one
INSERT INTO pending_servers (
    id,
    user_id,
    registration_token_hash,
    created_at,
    expires_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPendingServersByUserId :many
SELECT * FROM pending_servers
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetPendingServerByPendingServerId :one
SELECT * FROM pending_servers
WHERE id = $1;

-- name: DeletePendingServerByPendingServerId :exec
DELETE FROM pending_servers
WHERE id = $1;

-- name: CreateServerByUserId :exec
INSERT INTO servers (
    id,
    name,
    user_id,
    token_hash,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetServersByUserId :many
SELECT * FROM servers
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetServerByServerId :one
SELECT * FROM servers
WHERE id = $1;

-- name: UpdateServerByServerId :one
UPDATE servers
SET
    name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteServerByServerId :exec
DELETE FROM servers
WHERE id = $1;
