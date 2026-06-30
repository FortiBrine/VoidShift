-- name: CreateNetwork :one
INSERT INTO networks (name, address, listen_port, public_key, private_key)
VALUES (?, ?, ?, ?, ?)
RETURNING id, name, address, listen_port, public_key, private_key;

-- name: GetNetwork :one
SELECT id, name, address, listen_port, public_key, private_key
FROM networks WHERE id = ?;

-- name: GetNetworks :many
SELECT id, name, address, listen_port, public_key, private_key
FROM networks;

-- name: DeleteNetwork :exec
DELETE FROM networks WHERE id = ?;
