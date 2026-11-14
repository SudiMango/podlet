-- name: CreateEnvVarByUserIdAndRepoId :one
INSERT INTO env_vars (
    id,
    key,
    value_encrypted,
    user_id,
    repository_id
)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, repository_id, key)
DO UPDATE SET value_encrypted = EXCLUDED.value_encrypted
RETURNING id, key, user_id, repository_id, created_at;

-- name: GetEnvVarsByUserIdAndRepoId :many
SELECT id, key, user_id, repository_id, created_at FROM env_vars
WHERE user_id = $1
    AND repository_id = $2
ORDER BY key;

-- name: GetEnvVarsByUserIdAndRepoIdForAgent :many
SELECT key, value_encrypted FROM env_vars
WHERE user_id = $1
    AND repository_id = $2
ORDER BY key;

-- name: DeleteEnvVarByUserIdAndRepoId :execrows
DELETE FROM env_vars
WHERE user_id = $1
    AND repository_id = $2
    AND key = $3;

-- name: DeleteAllEnvVarsByUserIdAndRepoId :exec
DELETE FROM env_vars
WHERE user_id = $1
    AND repository_id = $2;
