-- name: CreateInstallationStateByUserId :one
INSERT INTO github_installation_states (
    id,
    user_id,
    created_at,
    expires_at
)
VALUES ($1, $2, $3, $4)
RETURNING *;