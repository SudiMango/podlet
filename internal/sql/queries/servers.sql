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

-- name: GetValidPendingServersByUserId :many
SELECT * FROM pending_servers
WHERE user_id = $1
    AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: GetValidPendingServerByUserIdAndPendingServerId :one
SELECT * FROM pending_servers
WHERE user_id = $1
    AND id = $2
    AND expires_at > NOW();

-- name: GetValidPendingServerRegistrationTokenHash :one
SELECT * FROM pending_servers
WHERE registration_token_hash = $1
    AND expires_at > NOW()
FOR UPDATE;

-- name: DeletePendingServerByUserIdAndPendingServerId :execrows
DELETE FROM pending_servers
WHERE user_id = $1
    AND id = $2;

-- name: CreateServerByUserId :one
INSERT INTO servers (
    id,
    name,
    user_id,
    token_hash,
    created_at,
    updated_at,
    last_seen_at
)
VALUES ($1, $2, $3, $4, $5, $6, NULL)
RETURNING *;

-- name: GetServersByUserId :many
SELECT * FROM servers
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetServerByUserIdAndServerId :one
SELECT * FROM servers
WHERE user_id = $1
    AND id = $2;

-- name: GetServerByTokenHash :one
SELECT * FROM servers
WHERE token_hash = $1;

-- name: UpdateServerByUserIdAndServerId :one
UPDATE servers
SET
    name = $3,
    updated_at = NOW()
WHERE user_id = $1
    AND id = $2
RETURNING *;

-- name: UpdateServerLastSeenAtByUserIdAndServerId :execrows
UPDATE servers
SET
    last_seen_at = $3
WHERE user_id = $1
    AND id = $2;

-- name: DeleteServerByUserIdAndServerId :execrows
DELETE FROM servers
WHERE user_id = $1
    AND id = $2;
