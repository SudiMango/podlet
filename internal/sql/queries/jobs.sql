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
    job_status = $3,
    error_message = CASE
        WHEN $3 = 'failed' THEN $4
        ELSE NULL
    END,
    updated_at = now()
WHERE id = $1
  AND server_id = $2
  AND (
    (job_status = 'pending' AND $3 = 'ongoing')
    OR (job_status = 'ongoing' AND $3 IN ('success', 'failed'))
  );
