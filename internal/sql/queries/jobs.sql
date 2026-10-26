-- name: CreatePendingJobForServer :one
INSERT INTO server_agent_jobs (
    id,
    user_id,
    server_id,
    job_type,
    repo_name,
    repository_id,
    installation_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetOldestJobByServerId :one
SELECT *
FROM server_agent_jobs
WHERE
    server_id = $1
    AND job_status = 'pending'
ORDER BY created_at ASC
LIMIT 1;