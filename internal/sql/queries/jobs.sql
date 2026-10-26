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

-- name: ReportJobByServerIdAndJobId :execrows
UPDATE server_agent_jobs
SET
    job_status = sqlc.arg(job_status)::server_job_status,
    error_message = CASE
        WHEN sqlc.arg(job_status)::server_job_status = 'failed' THEN sqlc.narg(error_message)
        ELSE NULL::TEXT
    END,
    updated_at = now()
WHERE id = sqlc.arg(id)
  AND server_id = sqlc.arg(server_id)
  AND (
    (job_status = 'pending' AND sqlc.arg(job_status)::server_job_status = 'ongoing')
    OR (job_status = 'ongoing' AND sqlc.arg(job_status)::server_job_status IN ('success', 'failed'))
  );
