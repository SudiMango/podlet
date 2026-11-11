/*
    Tunnel stuff
*/

-- name: CreateTunnel :one
INSERT INTO tunnels (
    id,
    name,
    user_id,
    server_id,
    created_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTunnelsByUserId :many
SELECT * FROM tunnels
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTunnelByUserIdAndServerId :one
SELECT * FROM tunnels
WHERE user_id = $1
    AND server_id = $2;

-- name: GetTunnelByIdAndUserId :one
SELECT * FROM tunnels
WHERE id = $1
    AND user_id = $2;

/*
    DNS stuff
*/

-- name: CreateDnsRecord :one
INSERT INTO dns_records (
    id,
    name,
    service,
    user_id,
    server_id,
    tunnel_id,
    repository_id,
    created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDnsRecordByServerIdAndRepositoryId :one
SELECT * FROM dns_records
WHERE user_id = $1
    AND server_id = $2
    AND repository_id = $3;

-- name: UpdateDnsRecordService :one
UPDATE dns_records
SET 
    service = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: GetDnsRecordsByUserId :many
SELECT * FROM dns_records
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetDnsRecordsByUserIdAndTunnelId :many
SELECT * FROM dns_records
WHERE user_id = $1
    AND tunnel_id = $2
ORDER BY created_at DESC;

-- name: GetDnsRecordByIdAndUserId :one
SELECT * FROM dns_records
WHERE id = $1
    AND user_id = $2;
