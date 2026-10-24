-- name: GetPodletProfileByUserId :one
SELECT
    u.id AS user_id,
    u.first_name,
    u.last_name,
    u.email,
    u.created_at,
    u.updated_at
FROM users u
WHERE u.id = $1;
