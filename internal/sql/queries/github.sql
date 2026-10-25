-- name: CreateInstallationStateByUserId :one
INSERT INTO github_installation_states (
    id,
    user_id,
    created_at,
    expires_at
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetStateByStateId :one
SELECT * FROM github_installation_states
WHERE id = $1;

-- name: UpsertGithubInstallation :one
INSERT INTO github_installations (
    id,
    user_id,
    account_id,
    account_login,
    account_type,
    repository_selection
)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    account_id = EXCLUDED.account_id,
    account_login = EXCLUDED.account_login,
    account_type = EXCLUDED.account_type,
    repository_selection = EXCLUDED.repository_selection,
    updated_at = now()
RETURNING *;

-- name: GetGithubInstallationsByUserId :many
SELECT * FROM github_installations
WHERE user_id = $1
ORDER BY created_at;

-- name: DeleteGithubInstallationState :exec
DELETE FROM github_installation_states
WHERE id = $1;
