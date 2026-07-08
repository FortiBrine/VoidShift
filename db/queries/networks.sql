-- name: CreateNetwork :one
INSERT INTO networks (name, address, listen_port, public_key, private_key)
VALUES (?, ?, ?, ?, ?)
RETURNING id, name, address, listen_port, public_key, private_key, enabled;

-- name: GetNetwork :one
SELECT id, name, address, listen_port, public_key, private_key, enabled
FROM networks WHERE id = ?;

-- name: GetNetworks :many
SELECT id, name, address, listen_port, public_key, private_key, enabled
FROM networks;

-- name: GetEnabledNetworks :many
SELECT id, name, address, listen_port, public_key, private_key, enabled
FROM networks WHERE enabled != 0;

-- name: SetNetworkEnabled :exec
UPDATE networks SET enabled = ? WHERE id = ?;

-- name: DeleteNetwork :exec
DELETE FROM networks WHERE id = ?;

-- name: CountNetworks :one
SELECT COUNT(*) FROM networks;
